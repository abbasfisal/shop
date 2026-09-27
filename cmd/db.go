package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbWipeCmd)
	dbWipeCmd.Flags().BoolVar(&dbWipeForce, "force", false, "skip the confirmation prompt")
}

var dbWipeForce bool

// wipeKind is one category of schema objects, together with the SQL that lists
// the DROP statements for every object of that category.
type wipeKind struct {
	Kind  string
	Query string
}

// Objects are wiped in slice order. Views go first so their drop count is not
// eaten by table CASCADE, types go last because everything else may depend on them.
var wipeKinds = []wipeKind{
	{"view", `
		SELECT 'DROP VIEW IF EXISTS ' || quote_ident(n.nspname) || '.' || quote_ident(c.relname) || ' CASCADE'
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('v', 'm')
			AND n.nspname !~ '^pg_'
			AND n.nspname <> 'information_schema'
			AND NOT EXISTS (SELECT 1 FROM pg_depend d WHERE d.objid = c.oid AND d.deptype = 'e')
		ORDER BY c.relname`},
	{"table", `
		SELECT 'DROP TABLE IF EXISTS ' || quote_ident(n.nspname) || '.' || quote_ident(c.relname) || ' CASCADE'
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('r', 'p', 'f')
			AND n.nspname !~ '^pg_'
			AND n.nspname <> 'information_schema'
			AND NOT EXISTS (SELECT 1 FROM pg_depend d WHERE d.objid = c.oid AND d.deptype = 'e')
		ORDER BY c.relname`},
	{"sequence", `
		SELECT 'DROP SEQUENCE IF EXISTS ' || quote_ident(n.nspname) || '.' || quote_ident(c.relname) || ' CASCADE'
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind = 'S'
			AND n.nspname !~ '^pg_'
			AND n.nspname <> 'information_schema'
			AND NOT EXISTS (SELECT 1 FROM pg_depend d WHERE d.objid = c.oid AND d.deptype = 'e')
		ORDER BY c.relname`},
	{"routine", `
		SELECT 'DROP ROUTINE IF EXISTS ' || quote_ident(n.nspname) || '.' || quote_ident(p.proname)
			|| '(' || pg_get_function_identity_arguments(p.oid) || ') CASCADE'
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		WHERE p.prokind IN ('f', 'p')
			AND n.nspname !~ '^pg_'
			AND n.nspname <> 'information_schema'
			AND NOT EXISTS (SELECT 1 FROM pg_depend d WHERE d.objid = p.oid AND d.deptype = 'e')
		ORDER BY p.proname`},
	{"type", `
		SELECT 'DROP TYPE IF EXISTS ' || quote_ident(n.nspname) || '.' || quote_ident(t.typname) || ' CASCADE'
		FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE n.nspname !~ '^pg_'
			AND n.nspname <> 'information_schema'
			AND (
				t.typtype IN ('d', 'e', 'r', 'm')
				OR (t.typtype = 'c' AND EXISTS (
					SELECT 1 FROM pg_class c WHERE c.oid = t.typrelid AND c.relkind = 'c'
				))
			)
			AND NOT EXISTS (SELECT 1 FROM pg_depend d WHERE d.objid = t.oid AND d.deptype = 'e')
		ORDER BY t.typname`},
}

type wipeStatement struct {
	Kind string
	DDL  string
}

var dbWipeCmd = &cobra.Command{
	Use:   "db:wipe",
	Short: "Drop every object in the database (tables, views, sequences, routines, types)",
	Long: `Drop all objects of the connected database — tables, views, materialized
views, sequences, functions, procedures, domains, enums, ranges and composite
types — while keeping the database itself and any installed extensions.
Equivalent to "php artisan db:wipe". Re-run "migrate" afterwards to rebuild.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		db := openMigrateDB()
		defer db.Close()

		dbName := currentDBName(db)
		confirmWipe(dbName)

		statements, err := collectWipeStatements(db)
		if err != nil {
			log.Fatal("-- db:wipe catalog query error: ", err)
		}
		if len(statements) == 0 {
			fmt.Printf("[db:wipe] database %q has no objects, nothing to drop\n", dbName)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("-- db:wipe begin transaction error: ", err)
		}
		for _, stmt := range statements {
			if _, err := tx.Exec(stmt.DDL); err != nil {
				_ = tx.Rollback()
				log.Fatalf("-- db:wipe failed on %s: %v\n-- statement: %s\n", stmt.Kind, err, stmt.DDL)
			}
		}
		if err := tx.Commit(); err != nil {
			log.Fatal("-- db:wipe commit error: ", err)
		}

		fmt.Printf("[db:wipe] dropped %d object(s) from %q (%s)\n",
			len(statements), dbName, summarizeKinds(statements))
		fmt.Println("[db:wipe] run `go run . migrate` to rebuild the schema")
	},
}

// currentDBName reports the connected database name for prompts and logs.
func currentDBName(db *sql.DB) string {
	var name string
	if err := db.QueryRow("SELECT current_database()").Scan(&name); err != nil {
		log.Fatal("-- db:wipe could not read current database: ", err)
	}
	return name
}

// confirmWipe blocks until the user types "yes" on a terminal. Non-interactive
// runs (CI, pipes) proceed, matching "php artisan db:wipe".
func confirmWipe(dbName string) {
	if dbWipeForce {
		return
	}
	if !isInteractiveTerminal() {
		fmt.Printf("[db:wipe] non-interactive session: wiping %q without confirmation\n", dbName)
		return
	}

	fmt.Printf("[db:wipe] this permanently deletes ALL objects in database %q. Type \"yes\" to continue: ", dbName)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("\n[db:wipe] aborted")
		os.Exit(1)
	}
	if strings.ToLower(strings.TrimSpace(line)) != "yes" {
		fmt.Println("[db:wipe] aborted")
		os.Exit(1)
	}
}

// collectWipeStatements asks the catalogs for every droppable object, in wipe order.
func collectWipeStatements(db *sql.DB) ([]wipeStatement, error) {
	var statements []wipeStatement
	for _, kind := range wipeKinds {
		rows, err := db.Query(kind.Query)
		if err != nil {
			return nil, fmt.Errorf("%s list: %w", kind.Kind, err)
		}
		for rows.Next() {
			var ddl string
			if err := rows.Scan(&ddl); err != nil {
				rows.Close()
				return nil, fmt.Errorf("%s scan: %w", kind.Kind, err)
			}
			statements = append(statements, wipeStatement{Kind: kind.Kind, DDL: ddl})
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, fmt.Errorf("%s rows: %w", kind.Kind, err)
		}
		rows.Close()
	}
	return statements, nil
}

func summarizeKinds(statements []wipeStatement) string {
	counts := make(map[string]int)
	for _, stmt := range statements {
		counts[stmt.Kind]++
	}
	parts := make([]string, 0, len(counts))
	for _, k := range wipeKinds {
		if counts[k.Kind] > 0 {
			parts = append(parts, fmt.Sprintf("%s: %d", k.Kind, counts[k.Kind]))
		}
	}
	return strings.Join(parts, ", ")
}
