package util

import "fmt"

// CategorySubtreeSQL returns a scalar sub-query (recursive CTE) that collects
// a category together with all of its descendants:
//
//	category_id IN (CategorySubtreeSQL(id)) AND ...
//
// Used by the storefront category listing and the slider catalog: a parent
// category like «مد و پوشاک» owns no products directly, only its children do.
// The id comes from the database as a uint, so inlining it is safe.
func CategorySubtreeSQL(categoryID uint) string {
	return fmt.Sprintf(`WITH RECURSIVE category_tree AS (
    SELECT id FROM categories WHERE id = %d
    UNION ALL
    SELECT c.id FROM categories c
             JOIN category_tree t ON c.parent_id = t.id
    WHERE c.deleted_at IS NULL
) SELECT id FROM category_tree`, categoryID)
}
