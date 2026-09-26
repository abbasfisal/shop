# Graph Report - shop  (2026-09-26)

## Corpus Check
- 549 files · ~2,791,151 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1612 nodes · 2747 edges · 152 communities (101 shown, 51 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 264 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6be4fd12`
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
- ProductService
- AdminLTE Page Skeleton Contract (head / navbar / sidebar / footer sub-templates + content-wrapper)
- ZapLogger
- BrandService
- HomeRepository
- CategoryRepository
- AdminHandler
- customer_order_response.go
- Storage
- Attribute
- Context
- NewHomeService
- GetAuthUser
- AdminHandler
- Duration
- web/customer_response.go
- SetErrors
- AttributeValueRepository
- product_read_model.go
- seeder.go
- .Payment
- New
- typesenceclient/typesence.go
- util.go
- AdminHandler
- Context
- Attribute
- ProductSliderService
- TestHandleError_InternalError
- AdminHandler
- AttributeValue
- BrandRepository
- send_email_task.go
- methods.go
- AdminHandler
- AuthenticateRepository
- ConvertGregorianToShamsi
- BannerService
- errors.go
- CustomerRepositoryInterface
- AuthenticateRepository
- SyncReadModel
- ProductRecommendation
- Context
- demo_shop.go
- NewRateLimiter
- AttributeService
- IsAdmin
- Context
- DashboardData
- Attribute
- OTP
- ProductAttribute
- ProductImages
- product_repository.go
- RefreshProductAggregates
- old.go
- testDB
- customer_payment_response.go
- AdminHandler
- Categories
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
- Context
- Brand
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
- CustomerRender
- Feature
- AGENTS.md
- ProductAttribute
- DB
- .PostLogin
- ToMenuResponse
- Pagination
- Model
- Customer
- Order
- Payment
- Attribute
- .IndexCustomer
- ProductSlider
- CreateBannerRequest
- CustomError
- Feature
- ProductImages
- Customer
- Order
- Attribute
- ProductVariant
- Engine
- Brand
- H
- T
- CreateBannerRequest
- UpdateProductRequest
- T
- Time
- VariantRow
- .Login
- Model
- Brand
- Category
- Product

## God Nodes (most connected - your core abstractions)
1. `HomeRepository` - 31 edges
2. `CustomError` - 29 edges
3. `ProductService` - 28 edges
4. `HandleError()` - 26 edges
5. `AdminHandler` - 25 edges
6. `HomeService` - 24 edges
7. `New()` - 21 edges
8. `AdminLTE Page Skeleton Contract (head / navbar / sidebar / footer sub-templates + content-wrapper)` - 21 edges
9. `SyncReadModel()` - 20 edges
10. `Render()` - 18 edges

## Surprising Connections (you probably didn't know these)
- `ValidateProductForm()` --calls--> `ValidProductStatus()`  [INFERRED]
  interfaces/http/requests/admin/product_form_validation.go → domain/entities/product.go
- `ValidateBannerForm()` --calls--> `ValidBannerLayout()`  [INFERRED]
  interfaces/http/requests/admin/banner_request.go → domain/entities/banner.go
- `ToBanner()` --calls--> `BannerLayoutLabel()`  [INFERRED]
  application/dto/admin/banner_response.go → domain/entities/banner.go
- `ToBanner()` --calls--> `BannerTypeLabel()`  [INFERRED]
  application/dto/admin/banner_response.go → domain/entities/banner.go
- `ToSlider()` --calls--> `SliderPositionLabel()`  [INFERRED]
  application/dto/admin/product_slider_response.go → domain/entities/product_slider.go

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

## Communities (152 total, 51 thin omitted)

### Community 0 - "CustomError"
Cohesion: 0.14
Nodes (12): CreateAttributeValueRequest, CreateCategoryRequest, AttributeValue, Context, Category, Context, Context, AttributeValueService (+4 more)

### Community 1 - "Product"
Cohesion: 0.08
Nodes (34): Product, ProductInventories, ProductInventory, Products, VariantAttributeValue, availableOf(), ProductVariant, Time (+26 more)

### Community 2 - "EventManager"
Cohesion: 0.06
Nodes (33): CancelJob, Context, Dependencies, registerSchedules(), RunScheduler(), Context, Dependencies, RunWorker() (+25 more)

### Community 3 - "GET"
Cohesion: 0.40
Nodes (5): Context, PublicHandler, ClearAll(), GET(), ValidateIRMobile()

### Community 4 - "Pagination"
Cohesion: 0.13
Nodes (12): UpdateOrderStatus, Context, NewOrderService(), OrderService, OrderServiceInterface, Pagination, buildLink(), buildLinkCurrent() (+4 more)

### Community 5 - "HomeService"
Cohesion: 0.09
Nodes (11): Category, Context, Customer, Order, Payment, HomeService, AddToCartRequest, CustomerProfileRequest (+3 more)

### Community 6 - "Site Header Partial"
Cohesion: 0.11
Nodes (45): Site 404 Not Found Page, Header-Only Error Page Shell, Site 500 Internal Error Page, Site Cart Page (define cart), Per-Product Cart Error Annotation (.ERRORS), Cart Line Quantity Mutation Forms, Customer Edit Profile Page, Account Navigation Active State (.ACTIVE) (+37 more)

### Community 7 - "admin_order_response.go"
Cohesion: 0.07
Nodes (48): Address, AdminOrder, AdminOrderItem, AdminOrderItems, AdminOrders, Customer, Customers, Order (+40 more)

### Community 8 - "Admin Sidebar Navigation Partial (sidebar)"
Cohesion: 0.13
Nodes (41): Admin AJAX Attribute Cascade (GET /admins/get-attributes/:id), Admin Authentication (mobile + password, POST /admins/login), Admin Route Navigation Contract (sidebar + navbar links), Admin Order Status Transition (states 0-11), Admin Page View-Model Key Contract (.TITLE/.MESSAGE/.ERRORS/.OLDS/.AUTH), Read-only Form Pattern (DISABLED inputs in a formless table), AdminLTE 3 Admin Theme Contract, Attribute / AttributeValue Administration (+33 more)

### Community 9 - "ProductService"
Cohesion: 0.07
Nodes (30): ProductListQuery, Context, NewPricingService(), Context, T, TestRefreshProductAggregates_DelegatesToRepository(), TestRefreshProductAggregates_PropagatesError(), Context (+22 more)

### Community 10 - "AdminLTE Page Skeleton Contract (head / navbar / sidebar / footer sub-templates + content-wrapper)"
Cohesion: 0.12
Nodes (35): admin_add_product_feature Template (Add Product Feature/صفات), Product Sub-Resource Navigation Hub (show-feature, show-gallery, add-attributes, add-inventory, product view/edit), admin_create_attribute Template (Create Attribute), Form State Contract (.OLDS repopulation, .ERRORS field errors, .MESSAGE flash), admin_create_attribute_values Template (Create Attribute Value), Attribute Hierarchy Contract (ATTRIBUTES.Data each carrying nested AttributeValues.Data; AJAX GET /admins/get-attributes/:id returns {Data:[{ID,Title}]}), Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form), create_banner Template (Banner Upload) (+27 more)

### Community 11 - "ZapLogger"
Cohesion: 0.09
Nodes (20): DB, openMigrateDB(), Close(), Connect(), env(), Get(), DB, Level (+12 more)

### Community 12 - "BrandService"
Cohesion: 0.18
Nodes (7): CreateBrandRequest, Brand, Context, NewBrandService(), BrandService, BrandServiceInterface, BrandRepositoryInterface

### Community 13 - "HomeRepository"
Cohesion: 0.06
Nodes (34): Customer, CustomerProfileRequest, CustomerVerifyRequest, Duration, HomeRepository, HomeRepositoryInterface, IncreaseCartItemQty, Context (+26 more)

### Community 14 - "CategoryRepository"
Cohesion: 0.10
Nodes (17): UpdateCategoryRequest, NewCategoryService(), CategoryRepository, CategoryServiceInterface, Category, Context, DB, NewCategoryRepository() (+9 more)

### Community 15 - "AdminHandler"
Cohesion: 0.16
Nodes (12): Banner, Banners, ToBanner(), ToBanners(), BannerLayoutLabel(), BannerTypeLabel(), Time, startOfDay() (+4 more)

### Community 16 - "customer_order_response.go"
Cohesion: 0.11
Nodes (29): CustomerOrderStatusMap(), Address, AttributeValue, Order, OrderItem, OrderItemAttribute, OrderItemAttributes, Time (+21 more)

### Community 17 - "Storage"
Cohesion: 0.09
Nodes (18): Time, ToCustomerSession(), Address, Cart, Model, Order, Customer, Model (+10 more)

### Community 18 - "Attribute"
Cohesion: 0.20
Nodes (11): Attribute, Attributes, ToAttribute(), ToAttributes(), AttributeRepository, AttributeRepositoryInterface, AttributeValues, CreateAttributeRequest (+3 more)

### Community 19 - "Context"
Cohesion: 0.33
Nodes (3): CreateProductInventoryRequest, Context, ProductRepository

### Community 20 - "NewHomeService"
Cohesion: 0.20
Nodes (6): Dependencies, NewHomeService(), HomeServiceInterface, HandlerFunc, LoadMenu(), HomeRepositoryInterface

### Community 21 - "GetAuthUser"
Cohesion: 0.15
Nodes (12): CheckCustomerSessionID(), HandlerFunc, CheckUserAuth(), HandlerFunc, CustomerMustLogin(), HandlerFunc, Auth(), CustomerAuth() (+4 more)

### Community 24 - "web/customer_response.go"
Cohesion: 0.26
Nodes (13): Address, Cart, CartItem, Address, toAddress(), toCart(), toCartItem(), toCartItems() (+5 more)

### Community 25 - "SetErrors"
Cohesion: 0.29
Nodes (6): Bundle, Context, PublicHandler, Context, Mapper(), SetErrors()

### Community 26 - "AttributeValueRepository"
Cohesion: 0.15
Nodes (10): UpdateAttributeValueRequest, NewAttributeValueService(), AttributeValueRepository, AttributeValueServiceInterface, Attribute, AttributeValue, Context, DB (+2 more)

### Community 27 - "product_read_model.go"
Cohesion: 0.30
Nodes (11): Time, B, C, F, FData, Img, ImgData, Inventory (+3 more)

### Community 28 - "seeder.go"
Cohesion: 0.13
Nodes (22): Brand, Category, Model, Product, Time, SliderPositionLabel(), SliderPositions(), ValidSliderPosition() (+14 more)

### Community 29 - ".Payment"
Cohesion: 0.25
Nodes (5): SendOTP(), SendSuccShop(), Context, PublicHandler, token

### Community 30 - "New"
Cohesion: 0.18
Nodes (13): New(), NewZarinpal(), Number, paymentRequestReqBody, paymentRequestResp, paymentVerificationReqBody, paymentVerificationResp, refreshAuthorityReqBody (+5 more)

### Community 31 - "typesenceclient/typesence.go"
Cohesion: 0.29
Nodes (10): boolPtr(), Connect(), CreateSchema(), GetTClient(), Client, RecreateSchema(), DeleteInTypesence(), Context (+2 more)

### Community 32 - "util.go"
Cohesion: 0.06
Nodes (33): AddToCartRequest, checkIDAndExistence(), AdminHandler, Context, AdminHandler, Context, Context, HomeServiceInterface (+25 more)

### Community 33 - "AdminHandler"
Cohesion: 0.16
Nodes (13): AttributeServiceInterface, AttributeValueServiceInterface, AuthenticateServiceInterface, BrandServiceInterface, CategoryServiceInterface, CustomerServiceInterface, DashboardService, AdminHandler (+5 more)

### Community 34 - "Context"
Cohesion: 0.24
Nodes (5): CreateProductFeatureRequest, UpdateProductFeatureRequest, Context, Feature, ProductRepository

### Community 37 - "ProductSliderService"
Cohesion: 0.06
Nodes (40): CreateProductSliderRequest, VariantRow, Context, Pagination, ProductSlider, NewProductSliderService(), ValidBannerLayout(), BannerDates() (+32 more)

### Community 38 - "TestHandleError_InternalError"
Cohesion: 0.60
Nodes (4): T, TestHandleError_InternalError(), TestHandleError_RecordNotFound(), TestNewAndError()

### Community 40 - "AttributeValue"
Cohesion: 0.90
Nodes (4): AttributeValue, AttributeValues, ToAttributeValue(), ToAttributeValues()

### Community 41 - "BrandRepository"
Cohesion: 0.27
Nodes (6): UpdateBrandRequest, BrandRepository, Brand, Context, DB, NewBrandRepository()

### Community 42 - "send_email_task.go"
Cohesion: 0.33
Nodes (7): Context, Dependencies, Task, NewSendEmailJob(), TaskSendEmail(), SendEmailJob, SendEmailPayload

### Community 43 - "methods.go"
Cohesion: 0.24
Nodes (8): Context, IsGuest(), Flash(), Context, Engine, Remove(), Set(), Start()

### Community 45 - "AuthenticateRepository"
Cohesion: 0.17
Nodes (9): NewAuthenticateService(), AuthenticateRepository, AuthenticateService, AuthenticateServiceInterface, Context, DB, User, NewAuthenticateRepository() (+1 more)

### Community 46 - "ConvertGregorianToShamsi"
Cohesion: 0.67
Nodes (3): ConvertGregorianToShamsi(), ConvertShamsiToGregorian(), Time

### Community 47 - "BannerService"
Cohesion: 0.05
Nodes (41): CreateBannerRequest, Banner, Context, NewBannerService(), BannerRepository, BannerService, Context, Dependencies (+33 more)

### Community 48 - "errors.go"
Cohesion: 0.27
Nodes (8): MyError, Add(), errorMessages(), Get(), GetErrorMessage(), NewMyError(), SetFromErrors(), ToString()

### Community 49 - "CustomerRepositoryInterface"
Cohesion: 0.16
Nodes (9): NewCustomerService(), CustomerRepository, CustomerService, CustomerServiceInterface, Context, Customer, DB, NewCustomerRepository() (+1 more)

### Community 50 - "AuthenticateRepository"
Cohesion: 0.38
Nodes (5): AuthenticateRepository, Context, Customer, DB, NewAuthenticateRepository()

### Community 51 - "SyncReadModel"
Cohesion: 0.24
Nodes (11): derefUint(), derefUintToInt64(), Context, DB, Feature, Product, ProductImages, SyncReadModel() (+3 more)

### Community 52 - "ProductRecommendation"
Cohesion: 0.40
Nodes (3): Time, ProductRecommendation, RecommendedProduct

### Community 53 - "Context"
Cohesion: 0.38
Nodes (3): Context, ProductRepository, parsedIDs()

### Community 54 - "demo_shop.go"
Cohesion: 0.27
Nodes (18): copyBannerAsset(), demoAttrValueID(), demoBanners(), demoBrandID(), demoCategoryID(), demoProducts(), demoSliderProductIDs(), demoSliders() (+10 more)

### Community 55 - "NewRateLimiter"
Cohesion: 0.33
Nodes (5): HandlerFunc, NewRateLimiter(), Limit, Limiter, RateLimiter

### Community 56 - "AttributeService"
Cohesion: 0.18
Nodes (7): CreateAttributeRequest, Attribute, Context, NewAttributeService(), AttributeService, AttributeServiceInterface, AttributeRepositoryInterface

### Community 59 - "DashboardData"
Cohesion: 0.16
Nodes (14): NewDashboardService(), DashboardRepository, DashboardService, Time, DashboardData, DashboardStats, LowStockProduct, NewUser (+6 more)

### Community 60 - "Attribute"
Cohesion: 0.40
Nodes (4): AttributeValue, DB, Model, Attribute

### Community 62 - "ProductAttribute"
Cohesion: 0.33
Nodes (5): Attribute, AttributeValue, Model, Product, ProductAttribute

### Community 63 - "ProductImages"
Cohesion: 0.36
Nodes (7): ImageProduct, ImageProducts, ToImageProduct(), ToImageProducts(), Model, Product, ProductImages

### Community 64 - "product_repository.go"
Cohesion: 0.12
Nodes (22): Attribute, Context, ProductRepository, createVariantLinks(), createVariantWithLinks(), escapeLike(), existingVariantCombos(), Context (+14 more)

### Community 65 - "RefreshProductAggregates"
Cohesion: 0.33
Nodes (4): Context, DB, ProductRepository, RefreshProductAggregates()

### Community 66 - "old.go"
Cohesion: 0.40
Nodes (4): Get(), Context, Set(), ToString()

### Community 67 - "testDB"
Cohesion: 0.50
Nodes (7): DB, T, TestAttributeAfterCreateGeneratesCode(), testDB(), TestRefreshProductAggregates(), TestSyncReadModelAndStorefrontGetProduct(), uintPtr()

### Community 68 - "customer_payment_response.go"
Cohesion: 0.60
Nodes (4): Time, PaymentStatusMapper(), ToPayment(), Payment

### Community 70 - "Categories"
Cohesion: 0.90
Nodes (4): Categories, Category, ToCategories(), ToCategory()

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

### Community 82 - "Brand"
Cohesion: 0.90
Nodes (4): Brand, Brands, ToBrand(), ToBrands()

### Community 88 - "HTTP Error Page Templates"
Cohesion: 1.00
Nodes (3): HTTP Error Page Templates, 404 Error Page (no {{define}} wrapper), 500 Error Page (templates/html/errors/500)

### Community 109 - "CustomerRender"
Cohesion: 0.07
Nodes (34): ProductSlider, ProductSliders, SliderPosition, SliderProduct, SliderPositions(), SlidersByPosition(), ToSlider(), ToSliderProduct() (+26 more)

### Community 115 - "Feature"
Cohesion: 0.90
Nodes (4): Feature, Features, ToFeature(), ToFeatures()

### Community 116 - "AGENTS.md"
Cohesion: 0.29
Nodes (5): Commands, Gotchas, graphify, Layout & dependency direction, Tests

### Community 117 - "ProductAttribute"
Cohesion: 0.90
Nodes (4): ProductAttribute, ProductAttributes, ToProductAttribute(), ToProductAttributes()

### Community 120 - "ToMenuResponse"
Cohesion: 0.60
Nodes (3): Category, ToMenuResponse(), CategoryResponse

### Community 151 - ".Login"
Cohesion: 0.29
Nodes (5): LoginRequest, User, ToUserResponse(), Context, User

## Ambiguous Edges - Review These
- `modules/admin/html/admin_create_category Template (Create Category)` → `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`  [AMBIGUOUS]
  templates/admin/admin_create_attribute_values.html · relation: conceptually_related_to
- `Admin Login Page (modules/admin/html/admin_login)` → `Admin Head Partial (head)`  [AMBIGUOUS]
  templates/admin/admin_login.html · relation: references

## Knowledge Gaps
- **37 isolated node(s):** `AdminHandler`, `Layout & dependency direction`, `Commands`, `Gotchas`, `Tests` (+32 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **51 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `modules/admin/html/admin_create_category Template (Create Category)` and `Orphan #parent-category AJAX Handler (dead copy-paste: change listener for a parent-category select that does not exist on this form)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Admin Login Page (modules/admin/html/admin_login)` and `Admin Head Partial (head)`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `New()` connect `New` to `CustomError`, `EventManager`, `HomeService`, `TestHandleError_InternalError`, `ProductService`, `ZapLogger`, `BrandService`, `ConvertGregorianToShamsi`, `Storage`, `.Login`?**
  _High betweenness centrality (0.219) - this node is a cross-community bridge._
- **Why does `CustomerRender()` connect `CustomerRender` to `Context`, `GET`, `.Payment`, `SetErrors`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `NewZarinpal()` connect `New` to `EventManager`, `.Payment`?**
  _High betweenness centrality (0.126) - this node is a cross-community bridge._
- **Are the 23 inferred relationships involving `HandleError()` (e.g. with `.Index()` and `.Show()`) actually correct?**
  _`HandleError()` has 23 INFERRED edges - model-reasoned connections that need verification._
- **What connects `AdminHandler`, `Layout & dependency direction`, `Commands` to the rest of the system?**
  _37 weakly-connected nodes found - possible documentation gaps or missing edges._