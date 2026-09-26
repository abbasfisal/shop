# Graph Report - shop  (2026-09-26)

## Corpus Check
- 535 files · ~2,761,192 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1322 nodes · 2428 edges · 126 communities (107 shown, 19 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 393 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7c69e494`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CustomError
- Product
- EventManager
- GET
- Pagination
- HomeService
- Site Header Partial
- admin_order_response.go
- Admin Sidebar Navigation Partial (sidebar)
- ProductRepository
- AdminLTE Page Skeleton Contract (head / navbar / sidebar / footer sub-templates + content-wrapper)
- ZapLogger
- BrandRepository
- HomeRepository
- CategoryRepository
- util.go
- customer_order_response.go
- Storage
- Seed
- SyncReadModel
- CustomerRender
- SetPublic
- Render
- .OrderPaidSuccessfully
- web/customer_response.go
- AdminHandler
- AttributeValueRepository
- product_read_model.go
- Add
- checkIDAndExistence
- New
- typesenceclient/typesence.go
- Error500
- SetAdminRoutes
- Context
- AttributeRepository
- customerWithGlobalData
- zarinpal.go
- postAddToCart
- AdminHandler
- Attribute
- DashboardData
- send_email_task.go
- banner.go
- AdminHandler
- AuthenticateRepository
- order_response.go
- NewAuthenticateService
- SetErrors
- NewCustomerService
- AuthenticateRepository
- setupRoutes
- RecommendedProduct
- Context
- Context
- NewRateLimiter
- NewAttributeService
- NewAttributeValueService
- NewCategoryService
- DashboardService
- Attribute
- .NewOtp
- ProductAttribute
- TestGetProduct_ReadModelContract
- Context
- RefreshProductAggregates
- old.go
- .IndexOrders
- customer_payment_response.go
- errors.go
- TestHandleError_InternalError
- Order
- Context
- HandleExampleTask
- HandleTaskSendWelcomeSMS
- AttributeValue
- Brand
- Cart
- Category
- OrderItem
- user_created_event.go
- .ShowTypeSenceForm
- ConvertGregorianToShamsi
- Address
- CartItem
- Feature
- Payment
- User
- HTTP Error Page Templates
- LoadHtml
- customer_authenticate_repository_interface.go
- customer_login_request.go
- verify_payment_query_string.go
- GetDep
- AGENTS.md
- RunScheduler
- CancelJob
- Auth
- RunWorker
- kavenegar.go
- .PostLogin
- .IndexCustomer
- IsAdmin
- IsGuest

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

## Communities (126 total, 19 thin omitted)

### Community 0 - "CustomError"
Cohesion: 0.06
Nodes (32): Attributes, Categories, Category, CreateAttributeRequest, CreateAttributeValueRequest, CreateCategoryRequest, CreateProductRequest, LoginRequest (+24 more)

### Community 1 - "Product"
Cohesion: 0.06
Nodes (43): Feature, Features, ImageProduct, ImageProducts, Product, ProductAttribute, ProductAttributes, ProductInventories (+35 more)

### Community 2 - "EventManager"
Cohesion: 0.19
Nodes (11): EventManager, EventManagerDep, eventName, ListenerFunc, Client, Context, DB, NewEventManager() (+3 more)

### Community 3 - "GET"
Cohesion: 0.24
Nodes (10): Context, PublicHandler, ClearAll(), Flash(), GET(), Context, Engine, Set() (+2 more)

### Community 4 - "Pagination"
Cohesion: 0.10
Nodes (17): UpdateOrderStatus, Context, NewOrderService(), Context, Customer, DB, Order, NewOrderRepository() (+9 more)

### Community 5 - "HomeService"
Cohesion: 0.06
Nodes (20): Time, ToCustomerSession(), Category, ToMenuResponse(), Context, Customer, Dependencies, Order (+12 more)

### Community 6 - "Site Header Partial"
Cohesion: 0.11
Nodes (45): Site 404 Not Found Page, Header-Only Error Page Shell, Site 500 Internal Error Page, Site Cart Page (define cart), Per-Product Cart Error Annotation (.ERRORS), Cart Line Quantity Mutation Forms, Customer Edit Profile Page, Account Navigation Active State (.ACTIVE) (+37 more)

### Community 7 - "admin_order_response.go"
Cohesion: 0.09
Nodes (33): Address, AdminOrder, AdminOrderItem, AdminOrderItems, AdminOrders, Customer, Customers, OrderDetail (+25 more)

### Community 8 - "Admin Sidebar Navigation Partial (sidebar)"
Cohesion: 0.13
Nodes (41): Admin AJAX Attribute Cascade (GET /admins/get-attributes/:id), Admin Authentication (mobile + password, POST /admins/login), Admin Route Navigation Contract (sidebar + navbar links), Admin Order Status Transition (states 0-11), Admin Page View-Model Key Contract (.TITLE/.MESSAGE/.ERRORS/.OLDS/.AUTH), Read-only Form Pattern (DISABLED inputs in a formless table), AdminLTE 3 Admin Theme Contract, Attribute / AttributeValue Administration (+33 more)

### Community 9 - "ProductRepository"
Cohesion: 0.08
Nodes (24): UpdateProductRequest, Context, NewPricingService(), Context, T, TestRefreshProductAggregates_DelegatesToRepository(), TestRefreshProductAggregates_PropagatesError(), NewProductService() (+16 more)

### Community 10 - "AdminLTE Page Skeleton Contract (head / navbar / sidebar / footer sub-templates + content-wrapper)"
Cohesion: 0.12
Nodes (35): admin_add_product_feature Template (Add Product Feature/صفات), Product Sub-Resource Navigation Hub (show-feature, show-gallery, add-attributes, add-inventory, product view/edit), admin_create_attribute Template (Create Attribute), Form State Contract (.OLDS repopulation, .ERRORS field errors, .MESSAGE flash), admin_create_attribute_values Template (Create Attribute Value), Attribute Hierarchy Contract (ATTRIBUTES.Data each carrying nested AttributeValues.Data; AJAX GET /admins/get-attributes/:id returns {Data:[{ID,Title}]}), Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form), create_banner Template (Banner Upload) (+27 more)

### Community 11 - "ZapLogger"
Cohesion: 0.11
Nodes (15): DB, openMigrateDB(), Level, Category, ExtraKey, Logger, SubCategory, ZapLogger (+7 more)

### Community 12 - "BrandRepository"
Cohesion: 0.11
Nodes (16): Brand, Brands, CreateBrandRequest, UpdateBrandRequest, ToBrand(), ToBrands(), Brand, Context (+8 more)

### Community 13 - "HomeRepository"
Cohesion: 0.18
Nodes (6): HomeRepository, Category, Context, Customer, Product, GetAuthUser()

### Community 14 - "CategoryRepository"
Cohesion: 0.14
Nodes (14): UpdateCategoryRequest, CategoryRepository, Category, Context, DB, NewCategoryRepository(), Delete(), Get() (+6 more)

### Community 15 - "util.go"
Cohesion: 0.13
Nodes (20): checkNationalCode(), ConvertIfNotNil(), ConvertSliceIfNotNil(), formatWithCommas(), GeneratePageNumbers(), T, HasSuffix(), PrettyJson() (+12 more)

### Community 16 - "customer_order_response.go"
Cohesion: 0.16
Nodes (22): CustomerOrderStatusMap(), Address, AttributeValue, Order, OrderItem, OrderItemAttribute, OrderItemAttributes, Time (+14 more)

### Community 17 - "Storage"
Cohesion: 0.11
Nodes (15): Address, Cart, Model, Order, Customer, Model, Time, Customer (+7 more)

### Community 18 - "Seed"
Cohesion: 0.16
Nodes (16): Close(), Connect(), env(), Get(), DB, fakeAttributeAndValues(), fakeBrands(), fakeCategories() (+8 more)

### Community 19 - "SyncReadModel"
Cohesion: 0.21
Nodes (9): CreateProductInventoryRequest, Context, ProductRepository, derefUintToInt64(), Context, DB, Feature, SyncReadModel() (+1 more)

### Community 20 - "CustomerRender"
Cohesion: 0.25
Nodes (7): HomeServiceInterface, Context, Dependencies, PublicHandler, NewPublicHandler(), CustomerRender(), GetProductStoragePath()

### Community 21 - "SetPublic"
Cohesion: 0.12
Nodes (11): CheckCustomerSessionID(), HandlerFunc, CheckUserAuth(), HandlerFunc, CustomerMustLogin(), HandlerFunc, HandlerFunc, LoadMenu() (+3 more)

### Community 22 - "Render"
Cohesion: 0.29
Nodes (5): AdminHandler, Context, AdminHandler, Context, Render()

### Community 23 - ".OrderPaidSuccessfully"
Cohesion: 0.19
Nodes (9): Client, Dependencies, Duration, Payment, NewHomeRepository(), releaseLocks(), retryWithBackoff(), Context (+1 more)

### Community 24 - "web/customer_response.go"
Cohesion: 0.26
Nodes (13): Address, Cart, CartItem, Address, toAddress(), toCart(), toCartItem(), toCartItems() (+5 more)

### Community 25 - "AdminHandler"
Cohesion: 0.17
Nodes (6): BrandServiceInterface, AdminHandler, Context, Dependencies, NewAdminHandler(), OrderServiceInterface

### Community 26 - "AttributeValueRepository"
Cohesion: 0.24
Nodes (7): UpdateAttributeValueRequest, AttributeValueRepository, Attribute, AttributeValue, Context, DB, NewAttributeRepository()

### Community 27 - "product_read_model.go"
Cohesion: 0.27
Nodes (12): Time, B, C, F, FData, Img, ImgData, Inventory (+4 more)

### Community 28 - "Add"
Cohesion: 0.32
Nodes (7): AdminHandler, Context, Add(), SetFromErrors(), Remove(), AllowImageExtensions(), GenerateFilename()

### Community 29 - "checkIDAndExistence"
Cohesion: 0.39
Nodes (4): checkIDAndExistence(), AdminHandler, Brand, Context

### Community 30 - "New"
Cohesion: 0.33
Nodes (5): New(), NewZarinpal(), Context, PublicHandler, Zarinpal

### Community 31 - "typesenceclient/typesence.go"
Cohesion: 0.29
Nodes (10): boolPtr(), Connect(), CreateSchema(), GetTClient(), Client, RecreateSchema(), DeleteInTypesence(), Context (+2 more)

### Community 32 - "Error500"
Cohesion: 0.43
Nodes (3): AdminHandler, Context, Error500()

### Community 33 - "SetAdminRoutes"
Cohesion: 0.13
Nodes (12): CreateBannerRequest, Context, NewBannerService(), BannerRepository, BannerService, Context, DB, NewBannerRepository() (+4 more)

### Community 34 - "Context"
Cohesion: 0.24
Nodes (5): CreateProductFeatureRequest, UpdateProductFeatureRequest, Context, Feature, ProductRepository

### Community 35 - "AttributeRepository"
Cohesion: 0.33
Nodes (5): AttributeRepository, Attribute, Context, DB, NewAttributeRepository()

### Community 36 - "customerWithGlobalData"
Cohesion: 0.31
Nodes (6): H, customerWithGlobalData(), Context, WithGlobalData(), StringToMap(), StringToUrlValues()

### Community 37 - "zarinpal.go"
Cohesion: 0.20
Nodes (10): Number, paymentRequestReqBody, paymentRequestResp, paymentVerificationReqBody, paymentVerificationResp, refreshAuthorityReqBody, refreshAuthorityResp, UnverifiedAuthority (+2 more)

### Community 38 - "postAddToCart"
Cohesion: 0.52
Nodes (6): T, postAddToCart(), TestAddToCart_ParsesNumericProductID(), TestAddToCart_RejectsLegacyObjectID(), ResponseRecorder, fakeHomeService

### Community 40 - "Attribute"
Cohesion: 0.47
Nodes (7): Attribute, AttributeValue, AttributeValues, ToAttribute(), ToAttributes(), ToAttributeValue(), ToAttributeValues()

### Community 41 - "DashboardData"
Cohesion: 0.42
Nodes (8): Time, DashboardData, DashboardStats, LowStockProduct, NewUser, PaymentReport, RecentOrder, StaticalReport

### Community 42 - "send_email_task.go"
Cohesion: 0.33
Nodes (7): Context, Dependencies, Task, NewSendEmailJob(), TaskSendEmail(), SendEmailJob, SendEmailPayload

### Community 43 - "banner.go"
Cohesion: 0.50
Nodes (3): Model, IsValidBannerType(), Banner

### Community 45 - "AuthenticateRepository"
Cohesion: 0.36
Nodes (5): AuthenticateRepository, Context, DB, User, NewAuthenticateRepository()

### Community 46 - "order_response.go"
Cohesion: 0.23
Nodes (15): Order, OrderItem, OrderItems, Orders, Payment, Payment, Time, OrderStatusMap() (+7 more)

### Community 47 - "NewAuthenticateService"
Cohesion: 0.33
Nodes (4): NewAuthenticateService(), AuthenticateService, AuthenticateServiceInterface, AuthenticateRepositoryInterface

### Community 48 - "SetErrors"
Cohesion: 0.29
Nodes (6): Bundle, Context, PublicHandler, Context, Mapper(), SetErrors()

### Community 49 - "NewCustomerService"
Cohesion: 0.16
Nodes (9): NewCustomerService(), CustomerRepository, CustomerService, CustomerServiceInterface, Context, Customer, DB, NewCustomerRepository() (+1 more)

### Community 50 - "AuthenticateRepository"
Cohesion: 0.38
Nodes (5): AuthenticateRepository, Context, Customer, DB, NewAuthenticateRepository()

### Community 51 - "setupRoutes"
Cohesion: 0.39
Nodes (7): Context, Dependencies, Engine, RunHttpServer(), RunServe(), setupRoutes(), setupSessions()

### Community 52 - "RecommendedProduct"
Cohesion: 0.33
Nodes (3): Time, ProductRecommendation, RecommendedProduct

### Community 53 - "Context"
Cohesion: 0.38
Nodes (3): Context, ProductRepository, parsedIDs()

### Community 55 - "NewRateLimiter"
Cohesion: 0.33
Nodes (5): HandlerFunc, NewRateLimiter(), Limit, Limiter, RateLimiter

### Community 56 - "NewAttributeService"
Cohesion: 0.33
Nodes (3): NewAttributeService(), AttributeServiceInterface, AttributeRepositoryInterface

### Community 57 - "NewAttributeValueService"
Cohesion: 0.33
Nodes (3): NewAttributeValueService(), AttributeValueServiceInterface, AttributeValueRepositoryInterface

### Community 58 - "NewCategoryService"
Cohesion: 0.33
Nodes (3): NewCategoryService(), CategoryServiceInterface, CategoryRepositoryInterface

### Community 59 - "DashboardService"
Cohesion: 0.24
Nodes (6): NewDashboardService(), DashboardRepository, DashboardService, DB, NewDashboardRepository(), DashboardRepositoryInterface

### Community 60 - "Attribute"
Cohesion: 0.40
Nodes (4): AttributeValue, DB, Model, Attribute

### Community 61 - ".NewOtp"
Cohesion: 0.33
Nodes (3): Model, OTP, Order

### Community 62 - "ProductAttribute"
Cohesion: 0.33
Nodes (5): Attribute, AttributeValue, Model, Product, ProductAttribute

### Community 63 - "TestGetProduct_ReadModelContract"
Cohesion: 0.47
Nodes (5): Context, T, keys(), nilGinCtx(), TestGetProduct_ReadModelContract()

### Community 64 - "Context"
Cohesion: 0.47
Nodes (3): Attribute, Context, ProductRepository

### Community 65 - "RefreshProductAggregates"
Cohesion: 0.33
Nodes (4): Context, DB, ProductRepository, RefreshProductAggregates()

### Community 66 - "old.go"
Cohesion: 0.40
Nodes (4): Get(), Context, Set(), ToString()

### Community 68 - "customer_payment_response.go"
Cohesion: 0.60
Nodes (4): Time, PaymentStatusMapper(), ToPayment(), Payment

### Community 69 - "errors.go"
Cohesion: 0.31
Nodes (6): MyError, errorMessages(), Get(), GetErrorMessage(), NewMyError(), ToString()

### Community 70 - "TestHandleError_InternalError"
Cohesion: 0.60
Nodes (4): T, TestHandleError_InternalError(), TestHandleError_RecordNotFound(), TestNewAndError()

### Community 71 - "Order"
Cohesion: 0.40
Nodes (4): Model, OrderItem, Payment, Order

### Community 73 - "HandleExampleTask"
Cohesion: 0.50
Nodes (4): Context, Task, HandleExampleTask(), TaskExample()

### Community 74 - "HandleTaskSendWelcomeSMS"
Cohesion: 0.50
Nodes (4): Context, Task, HandleTaskSendWelcomeSMS(), TaskSendWelcomeSMS()

### Community 75 - "AttributeValue"
Cohesion: 0.50
Nodes (3): Attribute, Model, AttributeValue

### Community 76 - "Brand"
Cohesion: 0.50
Nodes (3): Model, Product, Brand

### Community 77 - "Cart"
Cohesion: 0.50
Nodes (3): CartItem, Model, Cart

### Community 78 - "Category"
Cohesion: 0.50
Nodes (3): Model, Product, Category

### Community 79 - "OrderItem"
Cohesion: 0.50
Nodes (3): Model, Product, OrderItem

### Community 80 - "user_created_event.go"
Cohesion: 0.67
Nodes (3): Context, SendWelcomeNotification(), UserCreatedListener()

### Community 82 - "ConvertGregorianToShamsi"
Cohesion: 0.67
Nodes (3): ConvertGregorianToShamsi(), ConvertShamsiToGregorian(), Time

### Community 88 - "HTTP Error Page Templates"
Cohesion: 1.00
Nodes (3): HTTP Error Page Templates, 404 Error Page (no {{define}} wrapper), 500 Error Page (templates/html/errors/500)

### Community 115 - "GetDep"
Cohesion: 0.22
Nodes (7): SendEmailPayload, SyncReadModelEventPayload, GetDep(), Context, SendEmailListener(), Context, SyncReadModelListener()

### Community 116 - "AGENTS.md"
Cohesion: 0.29
Nodes (5): Commands, Gotchas, graphify, Layout & dependency direction, Tests

### Community 117 - "RunScheduler"
Cohesion: 0.38
Nodes (5): Context, Dependencies, registerSchedules(), RunScheduler(), Scheduler

### Community 118 - "CancelJob"
Cohesion: 0.47
Nodes (5): CancelJob, Dependencies, Task, NewCancelJob(), TaskCancelPendingOrders()

### Community 119 - "Auth"
Cohesion: 0.40
Nodes (5): Auth(), CustomerAuth(), Context, Customer, User

### Community 120 - "RunWorker"
Cohesion: 0.40
Nodes (3): Context, Dependencies, RunWorker()

### Community 121 - "kavenegar.go"
Cohesion: 0.40
Nodes (3): SendOTP(), SendSuccShop(), token

## Ambiguous Edges - Review These
- `modules/admin/html/admin_create_category Template (Create Category)` → `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`  [AMBIGUOUS]
  templates/admin/admin_create_attribute_values.html · relation: conceptually_related_to
- `Admin Login Page (modules/admin/html/admin_login)` → `Admin Head Partial (head)`  [AMBIGUOUS]
  templates/admin/admin_login.html · relation: references

## Knowledge Gaps
- **37 isolated node(s):** `Layout & dependency direction`, `Commands`, `Gotchas`, `Tests`, `graphify` (+32 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `modules/admin/html/admin_create_category Template (Create Category)` and `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Admin Login Page (modules/admin/html/admin_login)` and `Admin Head Partial (head)`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `New()` connect `New` to `CustomError`, `SetAdminRoutes`, `HomeService`, `TestHandleError_InternalError`, `postAddToCart`, `ProductRepository`, `ZapLogger`, `BrandRepository`, `HomeRepository`, `Storage`, `ConvertGregorianToShamsi`, `.OrderPaidSuccessfully`, `Add`, `.NewOtp`?**
  _High betweenness centrality (0.257) - this node is a cross-community bridge._
- **Why does `SyncReadModel()` connect `SyncReadModel` to `Context`, `Context`, `Context`, `ProductRepository`, `HomeRepository`, `Seed`, `GetDep`, `.OrderPaidSuccessfully`, `product_read_model.go`, `typesenceclient/typesence.go`, `TestGetProduct_ReadModelContract`?**
  _High betweenness centrality (0.134) - this node is a cross-community bridge._
- **Why does `SetAdminRoutes()` connect `SetAdminRoutes` to `Pagination`, `ProductRepository`, `BrandRepository`, `CategoryRepository`, `NewAuthenticateService`, `NewCustomerService`, `setupRoutes`, `NewRateLimiter`, `AdminHandler`, `NewAttributeService`, `NewAttributeValueService`, `NewCategoryService`, `DashboardService`, `New`?**
  _High betweenness centrality (0.129) - this node is a cross-community bridge._
- **Are the 45 inferred relationships involving `HandleError()` (e.g. with `.Index()` and `.Show()`) actually correct?**
  _`HandleError()` has 45 INFERRED edges - model-reasoned connections that need verification._
- **Are the 32 inferred relationships involving `Render()` (e.g. with `.CreateAttribute()` and `.CreateAttributeValues()`) actually correct?**
  _`Render()` has 32 INFERRED edges - model-reasoned connections that need verification._