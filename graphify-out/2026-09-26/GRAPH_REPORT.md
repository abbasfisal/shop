# Graph Report - .  (2026-09-26)

## Corpus Check
- Scoped corpus: 258 files, ~74,046 words (domain, application, infrastructure, interfaces, cmd, pkg, migrations, templates).

## Summary
- 1315 nodes · 2422 edges · 115 communities (100 shown, 15 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 393 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Admin DTOs and Requests
- Product Response DTOs
- CLI Serve and Scheduler
- SMS Provider and Error Types
- Admin Customer and Order DTOs
- Web Session and Menu DTOs
- Storefront Cart and Account Pages
- Admin Order DTOs
- Admin UI Contracts and Login
- Pricing Service and Aggregates
- Admin Create and Edit Forms
- Database Migrations Command
- Brand CRUD Vertical Slice
- Home Repository Implementation
- Category CRUD Vertical Slice
- Shared Utility Helpers
- Customer Order Web DTOs
- Customer and Session Entities
- Postgres Connection and Seeders
- Product Read Model Sync
- Public Web Handler Wiring
- Auth and Session Middleware
- Admin Handler Methods
- Order Payment and Background Jobs
- Customer Web Response Mapping
- Admin Inventory and Gallery Handlers
- Attribute Value CRUD Slice
- Product Read Model Struct
- Attribute Value Handler Methods
- Brand Handler Methods
- Payment Gateway Entry Point
- Typesense Search Client
- Category Handler Methods
- Banner CRUD Vertical Slice
- Product Feature CRUD Slice
- Attribute Repository Implementation
- Template Rendering Helpers
- Zarinpal Payment Wire Types
- Cart Handler Tests
- Attribute Handler Methods
- Attribute Response DTOs
- Dashboard Entity and Stats
- Email Background Task
- Banner Entity and Handlers
- Product Feature Handler Methods
- Admin Auth Repository
- Banner Repository and Routes
- Admin Auth Service
- Service Interfaces and Wiring
- Customer Service Layer
- Customer Session Repository
- Customer Repository Implementation
- Product Recommendation Entity
- Recommendation Repository
- Cart Handler Methods
- Rate Limiter Middleware
- Attribute Service Layer
- Attribute Value Service Layer
- Category Service Layer
- Dashboard Service Layer
- Attribute Domain Entity
- OTP Entity and Creation
- Product Attribute Link Entity
- Product Repository Integration Tests
- Product Attribute Repository
- Product Pricing Repository
- Legacy Config Store
- Order Handler Methods
- Customer Payment Web DTO
- Dashboard Repository Implementation
- Domain Error Tests
- Order Domain Entity
- Product Media Repository
- Example Background Task
- Welcome SMS Background Task
- Attribute Value Domain Entity
- Brand Domain Entity
- Cart Domain Entity
- Category Domain Entity
- Order Item Domain Entity
- User Created Event Listener
- Typesense Search Handlers
- Jalali Date Conversion
- Address Domain Entity
- Cart Item Domain Entity
- Feature Domain Entity
- Payment Domain Entity
- User Domain Entity
- HTTP Error Page Templates
- HTML Template Loading
- Customer Auth Repository Interface
- Customer Login Request
- Payment Verify Query Binding

## God Nodes (most connected - your core abstractions)
1. `CustomError` - 54 edges
2. `HandleError()` - 48 edges
3. `Render()` - 36 edges
4. `New()` - 34 edges
5. `HomeRepository` - 30 edges
6. `ProductService` - 27 edges
7. `SyncReadModel()` - 27 edges
8. `SetAdminRoutes()` - 25 edges
9. `HomeService` - 24 edges
10. `AdminHandler` - 24 edges

## Surprising Connections (you probably didn't know these)
- `ToCustomerOrder()` --calls--> `PrettyPrice()`  [INFERRED]
  application/dto/web/customer_order_response.go → pkg/util/util.go
- `toCartItem()` --calls--> `GetProductStoragePath()`  [INFERRED]
  application/dto/web/customer_response.go → pkg/util/util.go
- `SetAdminRoutes()` --calls--> `NewAttributeService()`  [INFERRED]
  interfaces/http/routes/admin_routes.go → application/usecases/attribute/attribute_service.go
- `SetAdminRoutes()` --calls--> `NewAttributeValueService()`  [INFERRED]
  interfaces/http/routes/admin_routes.go → application/usecases/attribute_value/attribute_value_service.go
- `SetAdminRoutes()` --calls--> `NewAuthenticateService()`  [INFERRED]
  interfaces/http/routes/admin_routes.go → application/usecases/auth/auth_service.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Product CRUD Workflow (create -> edit -> features -> gallery -> attributes -> inventory)** — templates_admin_admin_create_product_admin_create_product, templates_admin_admin_edit_product_admin_edit_product, templates_admin_admin_add_product_feature_admin_add_product_feature, templates_admin_admin_edit_product_feature_admin_edit_product_feature, templates_admin_admin_edit_gallery_product_edit_gallery_product, templates_admin_admin_create_product_att_values_admin_create_product_att_values, templates_admin_admin_create_product_inventory_inventory, templates_admin_admin_add_product_feature_product_subresource_nav [INFERRED 0.85]
- **Brand CRUD Triplet (list -> create -> edit)** — templates_admin_admin_index_brand_admin_index_brand, templates_admin_admin_create_brand_create_brand, templates_admin_admin_edit_brand_edit_brand [INFERRED 0.85]
- **Attribute Management Flow (attribute CRUD plus attribute-value CRUD over the nested attribute hierarchy)** — templates_admin_admin_index_attribute_admin_index_attribute, templates_admin_admin_create_attribute_admin_create_attribute, templates_admin_admin_index_attribute_value_admin_index_attribute_value, templates_admin_admin_create_attribute_values_admin_create_attribute_values, templates_admin_admin_create_attribute_values_attribute_hierarchy_contract, templates_admin_admin_create_attribute_form_state_contract [INFERRED 0.85]
- **AdminLTE Page Composition (head -> navbar -> sidebar -> content-wrapper -> footer)** — templates_layouts_head_head, templates_layouts_navbar_navbar, templates_layouts_sidebar_sidebar, templates_layouts_footer_footer, templates_layouts_starter_starter, templates_admin_admin_index_prodcut_modules_admin_html_admin_index_product, templates_admin_admin_show_attribute_admin_show_attribute, templates_admin_admin_show_attribute_values_edit_admin_show_attribute_values_edit, templates_admin_admin_show_attribute_values_index_admin_show_attribute_values_index, templates_admin_admin_show_brand_show_brand, templates_admin_admin_show_category_modules_admin_html_admin_show_category, templates_admin_admin_show_order_admin_show_order, templates_admin_admin_show_product_admin_show_product, templates_admin_admin_show_product_feature_admin_show_product_feature [EXTRACTED 1.00]
- **Admin View-Model Binding Contract (TITLE/MESSAGE/ERRORS/OLDS/AUTH + .Data envelope)** — admin_page_view_model, paginated_list_envelope, templates_admin_admin_index_prodcut_modules_admin_html_admin_index_product, templates_admin_admin_show_attribute_values_index_admin_show_attribute_values_index, templates_admin_admin_show_order_admin_show_order, templates_admin_admin_show_attribute_admin_show_attribute, templates_admin_admin_login_modules_admin_html_admin_login [INFERRED 0.95]
- **Session-Scoped Identity Rendering Across Admin and Customer Shells** — admin_authentication, customer_account_navigation, templates_admin_admin_login_modules_admin_html_admin_login, templates_layouts_sidebar_sidebar, templates_layouts_customer_sidebar_tmpl_customer_sidebar_tmpl, templates_admin_admin_show_attribute_admin_show_attribute [INFERRED 0.85]
- **Storefront Page Shell (header + mobile side menu + footers + go-to-top partials)** — templates_site_header_tmpl_header_tmpl, templates_site_mobile_side_menu_tmpl_mobile_side_menu_tmpl, templates_site_mobile_footer_menu_tmpl_mobile_footer_menu_tmpl, templates_site_main_footer_menu_tmpl_main_footer_menu_tmpl, templates_site_go_to_top_tmpl_go_to_top_tmpl, templates_site_home_home, templates_site_search_search, templates_site_single_product_single_product, templates_site_cart_cart, templates_site_customer_profile_customer_profile, templates_site_profile_orders_profile_orders, templates_site_customer_order_details_customer_order_details, templates_site_customer_edit_profile_customer_edit_profile, templates_site_shipping_shipping, templates_site_404_404, templates_site_500_500 [EXTRACTED 1.00]
- **Customer Authentication and Account Area (mobile login -> OTP verify -> profile/orders)** — templates_site_customer_login_customer_login, templates_site_customer_verify_phone_number_customer_verify_phone_number, templates_site_customer_login_passwordless_mobile_login, templates_site_customer_verify_phone_number_otp_code_entry, templates_site_customer_verify_phone_number_otp_resend_countdown, templates_site_customer_menu_tmpl_customer_menu_tmpl, templates_site_customer_menu_tmpl_account_navigation, templates_site_customer_profile_customer_profile, templates_site_customer_edit_profile_customer_edit_profile, templates_site_profile_orders_profile_orders, templates_site_customer_order_details_customer_order_details [INFERRED 0.85]
- **Browse-to-Cart-to-Checkout Flow (home/search -> product+variant -> cart lines -> shipping/payment -> order result)** — templates_site_home_home, templates_site_search_search, templates_site_single_product_single_product, templates_site_single_product_variant_inventory_picker, templates_site_single_product_add_to_cart_form, templates_site_single_product_already_in_cart_state, templates_site_cart_cart, templates_site_cart_cart_line_mutation_forms, templates_site_shipping_shipping, templates_site_shipping_checkout_summary_widget, templates_site_shopping_complete_buy_shopping_complete_buy, templates_site_shopping_complete_buy_checkout_steps, templates_site_shopping_no_complete_buy_shopping_no_complete_buy, templates_site_header_tmpl_mini_cart_widget [INFERRED 0.85]

## Communities (115 total, 15 thin omitted)

### Community 0 - "Admin DTOs and Requests"
Cohesion: 0.06
Nodes (33): Attributes, Categories, Category, CreateAttributeRequest, CreateAttributeValueRequest, CreateCategoryRequest, CreateProductRequest, LoginRequest (+25 more)

### Community 1 - "Product Response DTOs"
Cohesion: 0.07
Nodes (42): Feature, Features, ImageProduct, ImageProducts, Product, ProductAttribute, ProductAttributes, ProductInventories (+34 more)

### Community 2 - "CLI Serve and Scheduler"
Cohesion: 0.06
Nodes (38): CancelJob, Context, Dependencies, registerSchedules(), RunScheduler(), Context, Dependencies, Engine (+30 more)

### Community 3 - "SMS Provider and Error Types"
Cohesion: 0.06
Nodes (33): Bundle, MyError, SendOTP(), SendSuccShop(), Context, PublicHandler, Context, PublicHandler (+25 more)

### Community 4 - "Admin Customer and Order DTOs"
Cohesion: 0.07
Nodes (26): Customer, Customers, OrderDetail, UpdateOrderStatus, ToCustomer(), ToCustomers(), Customer, Order (+18 more)

### Community 5 - "Web Session and Menu DTOs"
Cohesion: 0.07
Nodes (18): Time, ToCustomerSession(), Category, ToMenuResponse(), Context, Customer, Dependencies, Order (+10 more)

### Community 6 - "Storefront Cart and Account Pages"
Cohesion: 0.11
Nodes (45): Site 404 Not Found Page, Header-Only Error Page Shell, Site 500 Internal Error Page, Site Cart Page (define cart), Per-Product Cart Error Annotation (.ERRORS), Cart Line Quantity Mutation Forms, Customer Edit Profile Page, Account Navigation Active State (.ACTIVE) (+37 more)

### Community 7 - "Admin Order DTOs"
Cohesion: 0.08
Nodes (40): Address, AdminOrder, AdminOrderItem, AdminOrderItems, AdminOrders, Order, OrderItem, OrderItemAttribute (+32 more)

### Community 8 - "Admin UI Contracts and Login"
Cohesion: 0.13
Nodes (41): Admin AJAX Attribute Cascade (GET /admins/get-attributes/:id), Admin Authentication (mobile + password, POST /admins/login), Admin Route Navigation Contract (sidebar + navbar links), Admin Order Status Transition (states 0-11), Admin Page View-Model Key Contract (.TITLE/.MESSAGE/.ERRORS/.OLDS/.AUTH), Read-only Form Pattern (DISABLED inputs in a formless table), AdminLTE 3 Admin Theme Contract, Attribute / AttributeValue Administration (+33 more)

### Community 9 - "Pricing Service and Aggregates"
Cohesion: 0.09
Nodes (23): UpdateProductRequest, Context, NewPricingService(), Context, T, TestRefreshProductAggregates_DelegatesToRepository(), TestRefreshProductAggregates_PropagatesError(), NewProductService() (+15 more)

### Community 10 - "Admin Create and Edit Forms"
Cohesion: 0.12
Nodes (35): admin_add_product_feature Template (Add Product Feature/صفات), Product Sub-Resource Navigation Hub (show-feature, show-gallery, add-attributes, add-inventory, product view/edit), admin_create_attribute Template (Create Attribute), Form State Contract (.OLDS repopulation, .ERRORS field errors, .MESSAGE flash), admin_create_attribute_values Template (Create Attribute Value), Attribute Hierarchy Contract (ATTRIBUTES.Data each carrying nested AttributeValues.Data; AJAX GET /admins/get-attributes/:id returns {Data:[{ID,Title}]}), Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form), create_banner Template (Banner Upload) (+27 more)

### Community 11 - "Database Migrations Command"
Cohesion: 0.11
Nodes (15): DB, openMigrateDB(), Level, Category, ExtraKey, Logger, SubCategory, ZapLogger (+7 more)

### Community 12 - "Brand CRUD Vertical Slice"
Cohesion: 0.11
Nodes (16): Brand, Brands, CreateBrandRequest, UpdateBrandRequest, ToBrand(), ToBrands(), Brand, Context (+8 more)

### Community 13 - "Home Repository Implementation"
Cohesion: 0.18
Nodes (6): HomeRepository, Category, Context, Customer, Product, GetAuthUser()

### Community 14 - "Category CRUD Vertical Slice"
Cohesion: 0.14
Nodes (14): UpdateCategoryRequest, CategoryRepository, Category, Context, DB, NewCategoryRepository(), Delete(), Get() (+6 more)

### Community 15 - "Shared Utility Helpers"
Cohesion: 0.13
Nodes (20): checkNationalCode(), ConvertIfNotNil(), ConvertSliceIfNotNil(), formatWithCommas(), GeneratePageNumbers(), T, HasSuffix(), PrettyJson() (+12 more)

### Community 16 - "Customer Order Web DTOs"
Cohesion: 0.16
Nodes (22): CustomerOrderStatusMap(), Address, AttributeValue, Order, OrderItem, OrderItemAttribute, OrderItemAttributes, Time (+14 more)

### Community 17 - "Customer and Session Entities"
Cohesion: 0.11
Nodes (15): Address, Cart, Model, Order, Customer, Model, Time, Customer (+7 more)

### Community 18 - "Postgres Connection and Seeders"
Cohesion: 0.16
Nodes (16): Close(), Connect(), env(), Get(), DB, fakeAttributeAndValues(), fakeBrands(), fakeCategories() (+8 more)

### Community 19 - "Product Read Model Sync"
Cohesion: 0.21
Nodes (9): CreateProductInventoryRequest, Context, ProductRepository, derefUintToInt64(), Context, DB, Feature, SyncReadModel() (+1 more)

### Community 20 - "Public Web Handler Wiring"
Cohesion: 0.25
Nodes (7): HomeServiceInterface, Context, Dependencies, PublicHandler, NewPublicHandler(), CustomerRender(), GetProductStoragePath()

### Community 21 - "Auth and Session Middleware"
Cohesion: 0.12
Nodes (11): CheckCustomerSessionID(), HandlerFunc, CheckUserAuth(), HandlerFunc, CustomerMustLogin(), HandlerFunc, HandlerFunc, LoadMenu() (+3 more)

### Community 22 - "Admin Handler Methods"
Cohesion: 0.23
Nodes (7): AdminHandler, Context, AdminHandler, Context, AdminHandler, Context, Render()

### Community 23 - "Order Payment and Background Jobs"
Cohesion: 0.19
Nodes (9): Client, Dependencies, Duration, Payment, NewHomeRepository(), releaseLocks(), retryWithBackoff(), Context (+1 more)

### Community 24 - "Customer Web Response Mapping"
Cohesion: 0.26
Nodes (13): Address, Cart, CartItem, Address, toAddress(), toCart(), toCartItem(), toCartItems() (+5 more)

### Community 26 - "Attribute Value CRUD Slice"
Cohesion: 0.24
Nodes (7): UpdateAttributeValueRequest, AttributeValueRepository, Attribute, AttributeValue, Context, DB, NewAttributeRepository()

### Community 27 - "Product Read Model Struct"
Cohesion: 0.27
Nodes (12): Time, B, C, F, FData, Img, ImgData, Inventory (+4 more)

### Community 28 - "Attribute Value Handler Methods"
Cohesion: 0.32
Nodes (4): AdminHandler, Context, Add(), SetFromErrors()

### Community 29 - "Brand Handler Methods"
Cohesion: 0.32
Nodes (5): checkIDAndExistence(), AdminHandler, Brand, Context, GenerateFilename()

### Community 30 - "Payment Gateway Entry Point"
Cohesion: 0.33
Nodes (5): New(), NewZarinpal(), Context, PublicHandler, Zarinpal

### Community 31 - "Typesense Search Client"
Cohesion: 0.29
Nodes (10): boolPtr(), Connect(), CreateSchema(), GetTClient(), Client, RecreateSchema(), DeleteInTypesence(), Context (+2 more)

### Community 32 - "Category Handler Methods"
Cohesion: 0.33
Nodes (4): AdminHandler, Context, Remove(), AllowImageExtensions()

### Community 33 - "Banner CRUD Vertical Slice"
Cohesion: 0.22
Nodes (6): CreateBannerRequest, Context, NewBannerService(), BannerService, Context, BannerRepositoryInterface

### Community 34 - "Product Feature CRUD Slice"
Cohesion: 0.24
Nodes (5): CreateProductFeatureRequest, UpdateProductFeatureRequest, Context, Feature, ProductRepository

### Community 35 - "Attribute Repository Implementation"
Cohesion: 0.33
Nodes (5): AttributeRepository, Attribute, Context, DB, NewAttributeRepository()

### Community 36 - "Template Rendering Helpers"
Cohesion: 0.29
Nodes (7): H, customerWithGlobalData(), Context, WithGlobalData(), StringToMap(), StringToUrlValues(), Flash()

### Community 37 - "Zarinpal Payment Wire Types"
Cohesion: 0.20
Nodes (10): Number, paymentRequestReqBody, paymentRequestResp, paymentVerificationReqBody, paymentVerificationResp, refreshAuthorityReqBody, refreshAuthorityResp, UnverifiedAuthority (+2 more)

### Community 38 - "Cart Handler Tests"
Cohesion: 0.27
Nodes (8): Context, T, postAddToCart(), TestAddToCart_ParsesNumericProductID(), TestAddToCart_RejectsLegacyObjectID(), ResponseRecorder, AddToCartRequest, fakeHomeService

### Community 39 - "Attribute Handler Methods"
Cohesion: 0.36
Nodes (3): AdminHandler, Context, Error500()

### Community 40 - "Attribute Response DTOs"
Cohesion: 0.47
Nodes (7): Attribute, AttributeValue, AttributeValues, ToAttribute(), ToAttributes(), ToAttributeValue(), ToAttributeValues()

### Community 41 - "Dashboard Entity and Stats"
Cohesion: 0.42
Nodes (8): Time, DashboardData, DashboardStats, LowStockProduct, NewUser, PaymentReport, RecentOrder, StaticalReport

### Community 42 - "Email Background Task"
Cohesion: 0.33
Nodes (7): Context, Dependencies, Task, NewSendEmailJob(), TaskSendEmail(), SendEmailJob, SendEmailPayload

### Community 43 - "Banner Entity and Handlers"
Cohesion: 0.29
Nodes (5): Model, IsValidBannerType(), Banner, AdminHandler, Context

### Community 45 - "Admin Auth Repository"
Cohesion: 0.36
Nodes (5): AuthenticateRepository, Context, DB, User, NewAuthenticateRepository()

### Community 46 - "Banner Repository and Routes"
Cohesion: 0.29
Nodes (6): BannerRepository, DB, NewBannerRepository(), Dependencies, Engine, SetAdminRoutes()

### Community 47 - "Admin Auth Service"
Cohesion: 0.33
Nodes (4): NewAuthenticateService(), AuthenticateService, AuthenticateServiceInterface, AuthenticateRepositoryInterface

### Community 48 - "Service Interfaces and Wiring"
Cohesion: 0.29
Nodes (4): BrandServiceInterface, Dependencies, NewAdminHandler(), ProductServiceInterface

### Community 49 - "Customer Service Layer"
Cohesion: 0.33
Nodes (4): NewCustomerService(), CustomerService, CustomerServiceInterface, CustomerRepositoryInterface

### Community 50 - "Customer Session Repository"
Cohesion: 0.38
Nodes (5): AuthenticateRepository, Context, Customer, DB, NewAuthenticateRepository()

### Community 51 - "Customer Repository Implementation"
Cohesion: 0.33
Nodes (5): CustomerRepository, Context, Customer, DB, NewCustomerRepository()

### Community 52 - "Product Recommendation Entity"
Cohesion: 0.33
Nodes (3): Time, ProductRecommendation, RecommendedProduct

### Community 53 - "Recommendation Repository"
Cohesion: 0.38
Nodes (3): Context, ProductRepository, parsedIDs()

### Community 55 - "Rate Limiter Middleware"
Cohesion: 0.33
Nodes (5): HandlerFunc, NewRateLimiter(), Limit, Limiter, RateLimiter

### Community 56 - "Attribute Service Layer"
Cohesion: 0.33
Nodes (3): NewAttributeService(), AttributeServiceInterface, AttributeRepositoryInterface

### Community 57 - "Attribute Value Service Layer"
Cohesion: 0.33
Nodes (3): NewAttributeValueService(), AttributeValueServiceInterface, AttributeValueRepositoryInterface

### Community 58 - "Category Service Layer"
Cohesion: 0.33
Nodes (3): NewCategoryService(), CategoryServiceInterface, CategoryRepositoryInterface

### Community 59 - "Dashboard Service Layer"
Cohesion: 0.47
Nodes (3): NewDashboardService(), DashboardService, DashboardRepositoryInterface

### Community 60 - "Attribute Domain Entity"
Cohesion: 0.40
Nodes (4): AttributeValue, DB, Model, Attribute

### Community 61 - "OTP Entity and Creation"
Cohesion: 0.33
Nodes (3): Model, OTP, Order

### Community 62 - "Product Attribute Link Entity"
Cohesion: 0.33
Nodes (5): Attribute, AttributeValue, Model, Product, ProductAttribute

### Community 63 - "Product Repository Integration Tests"
Cohesion: 0.47
Nodes (5): Context, T, keys(), nilGinCtx(), TestGetProduct_ReadModelContract()

### Community 64 - "Product Attribute Repository"
Cohesion: 0.47
Nodes (3): Attribute, Context, ProductRepository

### Community 65 - "Product Pricing Repository"
Cohesion: 0.33
Nodes (4): Context, DB, ProductRepository, RefreshProductAggregates()

### Community 66 - "Legacy Config Store"
Cohesion: 0.40
Nodes (4): Get(), Context, Set(), ToString()

### Community 68 - "Customer Payment Web DTO"
Cohesion: 0.60
Nodes (4): Time, PaymentStatusMapper(), ToPayment(), Payment

### Community 69 - "Dashboard Repository Implementation"
Cohesion: 0.50
Nodes (3): DashboardRepository, DB, NewDashboardRepository()

### Community 70 - "Domain Error Tests"
Cohesion: 0.60
Nodes (4): T, TestHandleError_InternalError(), TestHandleError_RecordNotFound(), TestNewAndError()

### Community 71 - "Order Domain Entity"
Cohesion: 0.40
Nodes (4): Model, OrderItem, Payment, Order

### Community 73 - "Example Background Task"
Cohesion: 0.50
Nodes (4): Context, Task, HandleExampleTask(), TaskExample()

### Community 74 - "Welcome SMS Background Task"
Cohesion: 0.50
Nodes (4): Context, Task, HandleTaskSendWelcomeSMS(), TaskSendWelcomeSMS()

### Community 75 - "Attribute Value Domain Entity"
Cohesion: 0.50
Nodes (3): Attribute, Model, AttributeValue

### Community 76 - "Brand Domain Entity"
Cohesion: 0.50
Nodes (3): Model, Product, Brand

### Community 77 - "Cart Domain Entity"
Cohesion: 0.50
Nodes (3): CartItem, Model, Cart

### Community 78 - "Category Domain Entity"
Cohesion: 0.50
Nodes (3): Model, Product, Category

### Community 79 - "Order Item Domain Entity"
Cohesion: 0.50
Nodes (3): Model, Product, OrderItem

### Community 80 - "User Created Event Listener"
Cohesion: 0.67
Nodes (3): Context, SendWelcomeNotification(), UserCreatedListener()

### Community 82 - "Jalali Date Conversion"
Cohesion: 0.67
Nodes (3): ConvertGregorianToShamsi(), ConvertShamsiToGregorian(), Time

### Community 88 - "HTTP Error Page Templates"
Cohesion: 1.00
Nodes (3): HTTP Error Page Templates, 404 Error Page (no {{define}} wrapper), 500 Error Page (templates/html/errors/500)

## Ambiguous Edges - Review These
- `modules/admin/html/admin_create_category Template (Create Category)` → `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`  [AMBIGUOUS]
  templates/admin/admin_create_attribute_values.html · relation: conceptually_related_to
- `Admin Login Page (modules/admin/html/admin_login)` → `Admin Head Partial (head)`  [AMBIGUOUS]
  templates/admin/admin_login.html · relation: references

## Knowledge Gaps
- **32 isolated node(s):** `CustomerAuthenticateRepositoryInterface`, `SendEmailPayload`, `SyncReadModelEventPayload`, `paymentRequestReqBody`, `paymentRequestResp` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `modules/admin/html/admin_create_category Template (Create Category)` and `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Admin Login Page (modules/admin/html/admin_login)` and `Admin Head Partial (head)`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `New()` connect `Payment Gateway Entry Point` to `Admin DTOs and Requests`, `Web Session and Menu DTOs`, `Domain Error Tests`, `Cart Handler Tests`, `Pricing Service and Aggregates`, `Database Migrations Command`, `Brand CRUD Vertical Slice`, `Home Repository Implementation`, `Banner Repository and Routes`, `Customer and Session Entities`, `Jalali Date Conversion`, `Brand Handler Methods`, `Order Payment and Background Jobs`, `OTP Entity and Creation`?**
  _High betweenness centrality (0.211) - this node is a cross-community bridge._
- **Why does `SyncReadModel()` connect `Product Read Model Sync` to `Product Attribute Repository`, `CLI Serve and Scheduler`, `Product Feature CRUD Slice`, `Product Media Repository`, `Pricing Service and Aggregates`, `Home Repository Implementation`, `Postgres Connection and Seeders`, `Order Payment and Background Jobs`, `Product Read Model Struct`, `Typesense Search Client`, `Product Repository Integration Tests`?**
  _High betweenness centrality (0.141) - this node is a cross-community bridge._
- **Why does `CustomError` connect `Admin DTOs and Requests` to `Web Session and Menu DTOs`, `Brand CRUD Vertical Slice`, `Order Payment and Background Jobs`, `OTP Entity and Creation`, `Payment Gateway Entry Point`?**
  _High betweenness centrality (0.105) - this node is a cross-community bridge._
- **Are the 45 inferred relationships involving `HandleError()` (e.g. with `.Index()` and `.Show()`) actually correct?**
  _`HandleError()` has 45 INFERRED edges - model-reasoned connections that need verification._
- **Are the 32 inferred relationships involving `Render()` (e.g. with `.CreateAttribute()` and `.CreateAttributeValues()`) actually correct?**
  _`Render()` has 32 INFERRED edges - model-reasoned connections that need verification._