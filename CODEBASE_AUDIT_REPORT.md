# ShoeGo E-Commerce Backend – Comprehensive Codebase Audit Report

---

## 1. Executive Summary

This report documents the findings of an exhaustive, read-only architectural, security, database, payment, and quality audit conducted on the **ShoeGo E-Commerce Backend** codebase. ShoeGo is an e-commerce platform built in Go (Golang) utilizing the Gin web framework, GORM ORM, PostgreSQL database, and Razorpay payment gateway integration.

The audit was executed under **strict read-only constraints**: no application source files, configuration files, database schemas, or dependencies were altered. 

### Key Audit Findings Overview

| Severity | Count | Primary Impact Areas |
| :--- | :---: | :--- |
| **CRITICAL** | **8** | Fatal server startup bug in `main.go`, public unauthenticated IDOR order disclosure, infinite coupon reuse (zero usage incrementing), wallet double-spending & broken atomicity, customer fund loss on wallet order cancellation, payment-stock mismatch with webhook 500 loops, hardcoded OAuth state (CSRF), and blocked user bypass. |
| **HIGH** | **13** | Plaintext OTP logging to stdout, cryptographically insecure PRNG OTP generation, inventory corruption on order cancellation, double refund on multi-step cancellations, non-expiring email change OTPs, full refund on partial returns, mismatched cookie names breaking cookie auth, zero rate limiting on auth endpoints, cart deletion of multi-variant products, sticky un-cleared coupons, 0777 image directory permissions, missing connection pooling, unmigrated domain models. |
| **MEDIUM** | **12** | Redundant config file loading on every request/helper call, missing Viper environment variable bindings, hardcoded `sslmode=disable`, misuse of 502 Bad Gateway for client validation errors, inconsistent non-envelope API responses, in-memory handoff map failing across instances, unchecked type assertions causing panics, unpaged sales report memory exhaustion, negative offset query crash, non-variant stock reporting 0 in cart, unindexed token blacklist bloat, dual JWT dependencies. |
| **LOW** | **7** | Ignored error returns (`_ = ...`), missing CORS middleware, missing standard security headers, inconsistent query parameters, REST route naming asymmetry, typographical errors in responses/logs, redundant PK model definitions. |
| **INFORMATIONAL** | **4** | 0% test coverage (zero test files in codebase), synchronous blocking SMTP email delivery, lack of structured logging (`slog`/`zap`), Swagger documentation drift. |
| **TOTAL FINDINGS** | **44** | |

### Production Readiness Verdict
**NOT PRODUCTION-READY.** The codebase exhibits multiple severe vulnerabilities and financial integrity flaws. In its current state, the server cannot start normally due to an inverted condition block in `main.go`, orders can be placed with infinite coupon discounts, wallet balances can be double-spent, cancelled wallet orders forfeit user funds without refund, and unauthenticated users can access customer order details. A phased remediation program must be enacted before deploying to any production or publicly accessible staging environment.

---

## 2. Project Overview

ShoeGo is a specialized footwear e-commerce RESTful API designed to manage customer shopping experiences and back-office administrative tasks.

### Core Capabilities
- **Customer Lifecycle:** User registration with email OTP verification, password login, Google OAuth 2.0 single sign-on, profile editing, multi-address management, password reset via OTP, and email address modification.
- **Catalog & Merchandising:** Category listing, product search with full-text ILIKE filtering, brand filtering, price-range constraints, multi-attribute sorting (price, name, new arrivals), and product variant management (sizes, variant SKUs, variant stock).
- **Cart & Wishlist:** Item additions with variant tracking, per-item quantity controls (max limit 5), wishlist toggling, and pre-checkout availability validation.
- **Promotions & Offers:** Percentage and fixed-amount coupon codes with minimum order thresholds and validity windows; product-level and category-level discount offers calculated dynamically.
- **Checkout & Payments:** Multi-address delivery selection, Cash on Delivery (COD), in-app Wallet payment, and Razorpay online gateway integration with HMAC-SHA256 signature verification and webhook processing.
- **Order Lifecycle:** Order tracking, invoice PDF generation, partial and full order cancellation, 7-day return request workflows, and wallet refunds.
- **Administration:** Admin authentication, customer block/unblock, inventory stock overrides, product/variant CRUD, category CRUD, promotional offer management, top-selling analytics (products, categories, brands), sales charts, and PDF sales report generation.

---

## 3. Technology Stack

| Component | Technology / Library | Version Specified in `go.mod` | Notes / Evaluation |
| :--- | :--- | :--- | :--- |
| **Language** | Go (Golang) | `go 1.26` (Runtime: `go1.26.0 darwin/arm64`) | Modern Go toolchain. |
| **HTTP Framework** | Gin Web Framework | `github.com/gin-gonic/gin v1.12.0` | High-performance HTTP multiplexer. |
| **ORM** | GORM | `gorm.io/gorm v1.31.1` | Used for associations, transactions, and queries. |
| **Database Driver** | PostgreSQL (pgx) | `gorm.io/driver/postgres v1.6.0` | Native Postgres driver via `jackc/pgx/v5`. |
| **Authentication** | JWT (JSON Web Tokens) | `golang-jwt/jwt/v5 v5.3.1` & `golang-jwt/jwt v3.2.2` | **Dual dependency issue**: v5 used for users; legacy v3 used for admin. |
| **Password Hashing**| Bcrypt | `golang.org/x/crypto v0.48.0` | Cost factor 10. |
| **OAuth 2.0** | Google OAuth | `golang.org/x/oauth2 v0.36.0` | Used for Google Sign-In. |
| **Payment Gateway** | Razorpay SDK | `github.com/razorpay/razorpay-go v1.4.1` | Order creation and signature validation. |
| **Configuration** | Viper | `github.com/spf13/viper v1.21.0` | Dotenv file loading and unmarshaling. |
| **Validation** | Go Playground Validator | `v10.30.1` & `v9.31.0+incompatible` | **Dual dependency issue**: v10 in config, deprecated v9 in handlers. |
| **Image Processing**| Imaging | `github.com/disintegration/imaging v1.6.2` | Image resizing and center cropping (800x800). |
| **PDF Generation** | Gofpdf | `github.com/jung-kurt/gofpdf v1.16.2` | Invoices and sales report generation. |
| **API Documentation**| Swaggo / Gin-Swagger | `github.com/swaggo/gin-swagger v1.6.1` | Swagger UI generation. |

---

## 4. Architecture and Codebase Review

### Architecture Pattern
The project implements a Layered Architecture (Handler-Usecase-Repository pattern):
$$\text{Client Request} \longrightarrow \text{Gin Router} \longrightarrow \text{Middleware} \longrightarrow \text{Handler} \longrightarrow \text{Usecase} \longrightarrow \text{Repository} \longrightarrow \text{PostgreSQL}$$

```
shoego-Ecommerce/
├── main.go                     # Application entry point & router bootstrap
├── config/                     # Configuration loader, OAuth, & Razorpay client initialization
├── database/                   # GORM connection & AutoMigrate declarations
├── domain/                     # Database entities and GORM model definitions
├── models/                     # DTOs, request/response structs, custom errors
├── handlers/                   # HTTP transport controllers (Gin context binding & serialization)
├── usecase/                    # Core business logic and workflow orchestration
├── repository/                 # Database queries, transactions, and persistence operations
├── middleware/                 # User authentication and admin authorization middlewares
├── helper/                     # Cryptography, JWT generation, image processing, PDF, & email utilities
├── response/                   # Standardized API response builder (`ClientResponse`)
├── router/                     # User and Admin route group declarations
├── template/                   # Static HTML templates (checkout handoff, login, payment failure)
└── docs/                       # Auto-generated Swagger documentation
```

### Architectural Strengths
1. **Clean Separation of Concerns:** Route registration is neatly segregated into `router/user.go` and `router/admin.go`. Business logic is largely separated from HTTP transport, facilitating future unit testability.
2. **Standardized Response Envelope:** The helper `response.ClientResponse(statusCode, message, data, error)` is used across most endpoints, providing consistency for frontend integrations.
3. **Database-Driven Token Invalidation:** Logout employs a server-side blacklist table (`blacklist_tokens`) to terminate active JWTs before natural expiry.

### Architectural Inconsistencies & Flaws
1. **Fatal Startup Logic Block in Entrypoint (`main.go`):**
   In `main.go` lines 57-72:
   ```go
   if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
       log.Fatal(err)
       router.Static("/images", "./images")
       router.Static("/static", "./template")
       routes.UserRoutes(router.Group("/"), db)
       routes.AdminRoutes(router.Group("/admin"), db)
       router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
       if err := router.Run("0.0.0.0:8080"); err != nil {
           log.Fatalf("Local host Errors %v", err)
       }
   }
   ```
   All route registrations, static asset declarations, and `router.Run` are nested **inside** the `if err != nil` error handling block of `SetTrustedProxies`. When `SetTrustedProxies` succeeds (`err == nil`), `main()` does not register any routes or execute `router.Run()`, causing the binary to exit immediately with status 0. If `SetTrustedProxies` fails, `log.Fatal(err)` terminates the process. As a result, the server literally cannot serve requests in its current form.
2. **Redundant Global Re-Reads of Configuration:**
   Instead of referencing the loaded `config.AppConfig` singleton, several helpers (`helper/auth.go:21`, `helper/admin.go:20,44`, `helper/helper.go:36`, `helper/sendotpemail.go:11`) re-read `.env` from disk and parse Viper configurations on **every single request** that performs authentication or sends an email.
3. **Admin Model & Table Disconnect:**
   `domain/admin.go` defines a `domain.Admin` struct, but `database/database.go` never migrates it. Meanwhile, `repository/admin.go` performs a raw SQL scan on the `users` table (`WHERE email = ? AND is_admin = true`). Admin identity is split between two conflicting concepts.
4. **Unmigrated Domain Model (`domain.Referral`):**
   `domain/offer.go:32` defines `type Referral struct`, yet it was omitted from the `AutoMigrate` sequence in `database/database.go`.
5. **No Context Propagation:**
   HTTP request contexts (`c.Request.Context()`) are not passed down to GORM database calls (`database.DB.WithContext(ctx)`). Client disconnections or HTTP timeouts cannot cancel running SQL queries, creating unnecessary load on PostgreSQL.

---

## 5. Complete API Route Audit

Every registered endpoint across `router/user.go`, `router/admin.go`, and `main.go` was inspected.

| # | Method | Endpoint | Feature | Status | Identified Issues | Severity |
| :-: | :--- | :--- | :--- | :---: | :--- | :---: |
| 1 | POST | `/signup` | User Registration | Flawed | Synchronous blocking SMTP email; plaintext OTP logged to server console. | HIGH |
| 2 | POST | `/verify-otp` | Registration Verification | Flawed | Ignores generated user token and returns nil data; predictable `math/rand` OTP verification. | HIGH |
| 3 | POST | `/resend-otp` | OTP Resend | Flawed | Plaintext OTP logged to server console; predictable PRNG seed; no IP rate limit. | HIGH |
| 4 | POST | `/login` | Password Authentication | Vulnerable | Returns 502 Bad Gateway for validation errors; sets cookie `"Authorization"` while middleware expects `"jwt_token"`; no rate limiting against brute force. | HIGH |
| 5 | POST | `/forgot-password` | Password Recovery | Flawed | Non-envelope response format (returns raw string); plaintext OTP logged; no rate limiting. | HIGH |
| 6 | POST | `/reset-password` | Password Reset | Working | Lacks rate limiting on OTP brute forcing. | MEDIUM |
| 7 | GET | `/auth/google/login` | Google OAuth Login | Vulnerable | Hardcoded OAuth `state="state"` parameter; completely vulnerable to OAuth CSRF attacks. | CRITICAL |
| 8 | GET | `/auth/google/callback` | Google OAuth Callback | Vulnerable | Does not verify OAuth `state`; does not check if user is `blocked`; enables blocked account bypass. | CRITICAL |
| 9 | GET | `/products` | Catalog Browsing | Flawed | Negative offset SQL error if `page=0`; Count query shares GORM session without isolation. | MEDIUM |
| 10 | GET | `/products/:id` | Product Details | Working | Public catalog inspection works as intended. | INFO |
| 11 | GET | `/categories` | Category Listing | Working | Returns listed categories alphabetically. | INFO |
| 12 | GET | `/order-success/:order_id` | Order Confirmation | Vulnerable | **Public IDOR**: Authorization check commented out; unauthenticated users can view full customer order details. | CRITICAL |
| 13 | GET | `/checkout/handoff` | Checkout Token Claim | Flawed | In-memory token storage fails in multi-pod/clustered deployments. | MEDIUM |
| 14 | POST | `/webhooks/razorpay` | Razorpay Webhook | Vulnerable | When inventory is depleted, rolls back payment status, leaves order pending, customer loses money, returns 500 causing endless retry loop. | CRITICAL |
| 15 | POST | `/logout` | User Logout | Flawed | Only inspects `Authorization` header; cannot log out users authenticated via cookie; does not clear cookie. | MEDIUM |
| 16 | GET | `/products/:id/validate`| Product Validation | Working | Verifies product listing and stock availability. | INFO |
| 17 | GET | `/referral` | User Referral Code | Working | Returns referral code and static link. | INFO |
| 18 | POST | `/cart` | Add to Cart | Working | Validates variant and listing status; removes item from wishlist. | INFO |
| 19 | GET | `/cart` | View Cart | Flawed | `Stock: item.Variant.Stock` returns 0 for all non-variant products. | MEDIUM |
| 20 | PATCH | `/cart/:product_id` | Update Cart Quantity | Flawed | Checks parent product stock instead of variant stock; affects all variants of a product. | HIGH |
| 21 | DELETE | `/cart/:product_id` | Remove Cart Item | Flawed | Deletes ALL variants of the specified product ID from cart instead of targeting a specific variant. | HIGH |
| 22 | GET | `/cart/validate` | Cart Checkout Validation | Working | Validates stock and availability before checkout. | INFO |
| 23 | POST | `/wishlist` | Add to Wishlist | Working | Prevents duplicates; validates product availability. | INFO |
| 24 | GET | `/wishlist` | View Wishlist | Working | Preloads products, categories, and images. | INFO |
| 25 | DELETE | `/wishlist/:product_id` | Remove from Wishlist | Working | Removes product from user wishlist. | INFO |
| 26 | POST | `/coupons/apply` | Apply Coupon to Cart | Vulnerable | Attaches coupon to `carts` table; never tracks `UsedCount`; unlimited coupon reuse possible. | CRITICAL |
| 27 | GET | `/coupons` | List Available Coupons | Working | Filters by subtotal and expiry date. | INFO |
| 28 | DELETE | `/coupons/remove` | Detach Coupon from Cart| Working | Clears `coupon_code` on active cart. | INFO |
| 29 | GET | `/checkout` | Checkout Summary Page | Working | Computes subtotal, tax (5%), shipping, and discounts. | INFO |
| 30 | POST | `/checkout/place-order` | Place COD Order | Flawed | Does not clear `cart.CouponCode`; coupon remains permanently applied to future orders. | HIGH |
| 31 | POST | `/checkout/wallet` | Place Wallet Order | Vulnerable | **Race Condition & Double Spend**: Non-atomic wallet deduction without `FOR UPDATE`; missing `domain.Payment` record forfeits 100% refund on cancellation. | CRITICAL |
| 32 | POST | `/checkout/handoff` | Create Handoff Link | Flawed | Generates link with in-memory map; fails across multiple instances. | MEDIUM |
| 33 | POST | `/checkout/razorpay` | Init Razorpay Payment | Flawed | Does not reserve inventory; duplicate clicks create orphan pending orders; no idempotency key. | HIGH |
| 34 | POST | `/checkout/razorpay/verify`| Verify Payment Signature| Vulnerable | If stock sold out during payment, rolls back transaction; customer charged with no order or refund. | CRITICAL |
| 35 | POST | `/checkout/razorpay/failure`| Record Payment Failure| Working | Validates payment status via Razorpay API before updating status to failed. | INFO |
| 36 | GET | `/orders` | User Order History | Working | Returns paginated orders belonging to user. | INFO |
| 37 | GET | `/orders/:order_id` | User Order Detail | Working | Properly scoped by `user_id`. | INFO |
| 38 | PATCH | `/orders/:order_id/cancel`| Full Order Cancellation | Vulnerable | Restores parent product stock for variant items; forfeits refund on wallet orders; double refunds if items were partially cancelled earlier. | CRITICAL |
| 39 | PATCH | `/orders/:order_id/items/:item_id/cancel`| Partial Cancellation| Flawed | Restores parent product stock AND variant stock, creating phantom parent inventory. | HIGH |
| 40 | PATCH | `/orders/:order_id/return`| Request Order Return | Working | Enforces 7-day post-delivery constraint. | INFO |
| 41 | PATCH | `/orders/:order_id/items/:item_id/return`| Request Item Return | Working | Sets item status to `return_requested`. | INFO |
| 42 | GET | `/orders/:order_id/invoice`| Download Invoice PDF | Working | Generates and streams PDF invoice. | INFO |
| 43 | GET | `/wallet` | View Wallet Balance | Working | Returns current user wallet balance. | INFO |
| 44 | GET | `/wallet/history` | View Wallet History | Working | Returns chronologically ordered wallet transactions. | INFO |
| 45 | GET | `/profile` | View User Profile | Working | Returns profile and associated address book. | INFO |
| 46 | PUT | `/profile/edit` | Edit Profile | Working | Unchecked type assertion on `user_id`. | LOW |
| 47 | PUT | `/profile/change-password`| Change Password | Working | Verifies old password hash with bcrypt. | INFO |
| 48 | POST | `/profile/request-email-change`| Request Email Change | Vulnerable | Generated OTP never expires (`ExpiresAt` is zero); does not require current password confirmation. | HIGH |
| 49 | POST | `/profile/verify-email-change`| Confirm Email Change | Vulnerable | Does not check OTP expiry; allows account email takeover. | HIGH |
| 50 | GET | `/address` | List User Addresses | Working | Returns addresses scoped to user. | INFO |
| 51 | POST | `/address` | Add Address | Working | Manages default address toggles. | INFO |
| 52 | PUT | `/address/:id` | Edit Address | Working | Scoped by `user_id` preventing IDOR. | INFO |
| 53 | DELETE | `/address/:id` | Delete Address | Working | Scoped by `user_id` preventing IDOR. | INFO |
| 54 | POST | `/admin/login` | Admin Authentication | Vulnerable | Returns 502 Bad Gateway for bad credentials; swagger doc specifies wrong path `/admin/admin`; no rate limiting. | HIGH |
| 55 | GET | `/admin/users` | List Users | Working | Paginated customer list with search. | INFO |
| 56 | PATCH | `/admin/block-user/:id` | Block User Account | Ineffective | Sets `blocked=true`, but blocked users can still log in via Google OAuth and keep using active JWTs. | HIGH |
| 57 | PATCH | `/admin/unblock-user/:id`| Unblock User Account | Working | Resets `blocked=false`. | INFO |
| 58 | POST | `/admin/products` | Create Product | Flawed | `0777` directory creation; no image upload size limit. | HIGH |
| 59 | PUT | `/admin/products/:id` | Update Product | Working | Supports partial multipart updates. | INFO |
| 60 | DELETE | `/admin/products/:id` | Soft Delete Product | Working | Soft deletes product using GORM. | INFO |
| 61 | GET | `/admin/products` | Admin Product List | Working | Paginated product list. | INFO |
| 62 | POST | `/admin/products/:id/variants`| Create Variant | Working | Associates sizes and stock to product. | INFO |
| 63 | GET | `/admin/products/:id/variants`| List Variants | Working | Returns variants for product. | INFO |
| 64 | PUT | `/admin/variants/:id` | Update Variant | Working | Updates size, SKU, and stock. | INFO |
| 65 | DELETE | `/admin/variants/:id` | Delete Variant | Working | Deletes product variant. | INFO |
| 66 | POST | `/admin/categories` | Create Category | Working | Prevents duplicate names case-insensitively. | INFO |
| 67 | PUT | `/admin/categories/:id`| Update Category | Working | Validates uniqueness excluding current ID. | INFO |
| 68 | DELETE | `/admin/categories/:id`| Delete Category | Working | Soft deletes category. | INFO |
| 69 | GET | `/admin/categories` | List Categories | Working | Paginated category list. | INFO |
| 70 | POST | `/admin/offers/product` | Create Product Offer | Working | Validates dates and discount bounds (1-90%). | INFO |
| 71 | GET | `/admin/offers/product` | List Product Offers | Working | Paginated offer list with active date calculation. | INFO |
| 72 | PUT | `/admin/offers/:id` | Update Product Offer | Flawed | Inconsistent route: `/offers/:id` updates product offer, but delete is `/offers/product/:id`. | LOW |
| 73 | DELETE | `/admin/offers/product/:id`| Delete Product Offer| Working | Deletes product offer. | INFO |
| 74 | POST | `/admin/category-offers`| Create Category Offer | Working | Validates dates and discount bounds (1-90%). | INFO |
| 75 | GET | `/admin/category-offers`| List Category Offers | Working | Paginated category offer list. | INFO |
| 76 | PUT | `/admin/category-offers/:id`| Update Category Offer| Working | Updates category offer. | INFO |
| 77 | DELETE | `/admin/category-offers/:id`| Delete Category Offer| Working | Deletes category offer. | INFO |
| 78 | GET | `/admin/orders` | Admin Order List | Working | Paginated orders with search and status filtering. | INFO |
| 79 | GET | `/admin/orders/:id` | Admin Order Detail | Working | Returns full order and customer breakdown. | INFO |
| 80 | PATCH | `/admin/orders/:id/status`| Update Order Status | Vulnerable | Approving a return on a `partially_returned` order refunds 100% of full order amount. | HIGH |
| 81 | GET | `/admin/inventory` | Inventory Management | Working | Stock status filtering (in stock, low stock, out of stock). | INFO |
| 82 | PATCH | `/admin/inventory/:id/stock`| Update Stock Level | Working | Overrides product stock. | INFO |
| 83 | GET | `/admin/dashboard/best-selling-products`| Analytics Products| Working | Aggregates top 10 sold products from delivered orders. | INFO |
| 84 | GET | `/admin/dashboard/best-selling-categories`| Analytics Categories| Working | Aggregates top 10 sold categories. | INFO |
| 85 | GET | `/admin/dashboard/best-selling-brands`| Analytics Brands| Working | Aggregates brand sales (returns brand IDs without names). | LOW |
| 86 | GET | `/admin/dashboard/sales-chart`| Sales Chart Analytics| Working | Aggregates sales grouped by daily, weekly, monthly, yearly. | INFO |
| 87 | POST | `/admin/coupons` | Create Coupon | Flawed | Does not validate `DiscountAmount > 0` or enforce percentage <= 100%. | MEDIUM |
| 88 | GET | `/admin/coupons` | List All Coupons | Working | Returns paginated coupons with active status. | INFO |
| 89 | PUT | `/admin/coupons/:id` | Update Coupon | Working | Updates coupon fields. | INFO |
| 90 | DELETE | `/admin/coupons/:id` | Delete Coupon | Working | Deletes coupon. | INFO |
| 91 | GET | `/admin/sales-report` | Sales Report JSON | Flawed | Unbounded order query loads all orders into memory; includes cancelled/refunded orders. | MEDIUM |
| 92 | GET | `/admin/sales-report/pdf`| Sales Report PDF | Flawed | Inconsistent query parameters (`period, from, to` vs `filter, start_date, end_date`). | LOW |
| 93 | POST | `/admin/logout` | Admin Logout | Working | Blacklists admin token string in database. | INFO |
| 94 | GET | `/images/*filepath` | Static Product Images | Flawed | Served from local filesystem with `0777` permissions; enclosed in broken startup block. | HIGH |
| 95 | GET | `/static/*filepath` | Static Templates | Flawed | Serves frontend checkout templates; enclosed in broken startup block. | HIGH |
| 96 | GET | `/swagger/*any` | Swagger Documentation | Flawed | Enclosed in broken startup block in `main.go`. | HIGH |

---

## 6. Feature-by-Feature Functional Audit

### 6.1 User Registration and Authentication
- **Workflow:** User submits signup data $\rightarrow$ Password hashed with bcrypt $\rightarrow$ Temporary record stored in `otp_verifications` $\rightarrow$ 6-digit OTP emailed $\rightarrow$ User submits OTP $\rightarrow$ Record verified and user created $\rightarrow$ Access (15m) and Refresh (90d) tokens generated.
- **Defects Identified:**
  1. `log.Println("OTP sent", otp)` exposes OTP in plain text to server logs.
  2. `helper.GenerateOTP()` uses `math/rand` with timestamp seed instead of `crypto/rand`.
  3. `handlers.VerifyOTP` discards the generated `models.TokenUser` with `_`, returning `nil` data to the client. Users cannot automatically proceed to authenticated state upon signup.
  4. User block status is not checked during Google OAuth or in JWT auth middleware.

### 6.2 Admin Authentication & User Management
- **Workflow:** Admin submits email/password $\rightarrow$ Repository queries `users` table where `is_admin = true` $\rightarrow$ Bcrypt compares password $\rightarrow$ Admin JWT generated with `KEY_FOR_ADMIN`.
- **Defects Identified:**
  1. `handlers.AdminLogin` returns `502 Bad Gateway` for validation and authentication errors.
  2. `repository.AdminLogin` runs raw SQL and scans into `domain.Admin`. If no record matches, GORM returns `nil` error and an empty struct. `usecase.AdminLogin` catches this only through `AdminDetail.Password == ""`, returning `"admin password is empty"` rather than `"admin email not found"`.
  3. Blocking a user (`BlockUser`) sets `blocked = true`, but this does not invalidate the user's active JWTs or prevent Google OAuth logins.

### 6.3 Product & Inventory Management
- **Workflow:** Admin uploads product with multipart form and minimum 3 images $\rightarrow$ Images decoded, resized, and saved to `./images` $\rightarrow$ Product variants added with individual size and stock.
- **Defects Identified:**
  1. Stock duality: `domain.Product` has `Stock`, and `domain.ProductVariant` has `Stock`. The relationship between parent product stock and variant stock is inconsistent across handlers.
  2. `usecase/product.go:AddProduct` saves images to disk before saving the product in the database. If `CreateProduct` fails, orphan images remain on disk.
  3. Image upload directories are created with `os.ModePerm` (`0777`).

### 6.4 Shopping Cart Management
- **Workflow:** User adds item with optional `variant_id` $\rightarrow$ Cart item created/incremented $\rightarrow$ Wishlist item automatically removed.
- **Defects Identified:**
  1. In `usecase/cart.go:GetCart` line 206, `Stock: item.Variant.Stock` is assigned directly. If the product does not have variants (`item.VariantID == nil`), the response reports `stock: 0`, falsely indicating the item is out of stock.
  2. `UpdateCartQuantity` and `RemoveCartItem` route on `product_id` instead of `item_id` or `(product_id, variant_id)`. If a user has two sizes of the same shoe in their cart, modifying or removing one destroys or modifies both.

### 6.5 Coupon Management
- **Workflow:** Admin creates coupon with code, discount type, amount, minimum spend, expiry date, and usage limit $\rightarrow$ Customer applies coupon $\rightarrow$ Discount calculated against cart subtotal $\rightarrow$ Code saved in `carts.coupon_code`.
- **Defects Identified:**
  1. **Infinite Coupon Abuse:** `coupon.UsedCount` is never incremented anywhere when orders are placed. `UsageLimit` is completely ineffective.
  2. **Sticky Coupon:** For COD and Wallet checkouts, `CreateOrderTransaction` deletes cart items but does not clear `coupon_code` on the `carts` row. The coupon automatically persists to the user's subsequent purchases.
  3. No per-user coupon tracking exists. The same user can apply the same promotional coupon indefinitely.

### 6.6 Checkout & Order Placement
- **Workflow:**
  - **COD:** `PlaceCODOrder` validates address $\rightarrow$ Creates order and order items $\rightarrow$ Deducts product/variant stock using atomic `Update("stock", gorm.Expr("stock - ?", item.Quantity))` $\rightarrow$ Clears cart items.
  - **Wallet:** `PlaceWalletOrder` checks balance $\rightarrow$ Creates order $\rightarrow$ Deducts wallet balance $\rightarrow$ Records wallet transaction.
  - **Razorpay:** `CreateRazorpayOrder` creates pending order in DB $\rightarrow$ Creates Razorpay gateway order $\rightarrow$ Returns order ID and key to client.
- **Defects Identified:**
  1. **Wallet Double-Spending:** `PlaceWalletOrder` does not lock the wallet record (`FOR UPDATE`). Concurrent requests from the same user can place multiple orders against the same balance.
  2. **Broken Wallet Atomicity:** Order creation and wallet balance updates run in separate database transactions. If the balance update fails, the order is already marked as paid.
  3. **Permanent Forfeiture of Wallet Funds on Cancellation:** `PlaceWalletOrder` fails to insert a record into the `payments` table. When the user cancels the order, `CancelOrder` looks for a payment record, finds none, and skips the refund.
  4. **Razorpay Inventory Race Condition:** `CreateRazorpayOrder` does not reserve or lock stock. If items sell out before payment verification, `ReduceStockAfterPaymentTx` rolls back, the order remains in "pending", and the customer's money is captured without fulfillment or automated refund.

### 6.7 Order Management, Returns & Refunds
- **Workflow:** User can cancel before delivery $\rightarrow$ Refund issued to wallet if paid $\rightarrow$ User can request return within 7 days of delivery $\rightarrow$ Admin approves return.
- **Defects Identified:**
  1. **Stock Corruption on Cancellation:** `CancelOrderItemTransaction` restores parent `Product.Stock` AND `ProductVariant.Stock`, artificially inflating parent stock.
  2. **Double Refund:** Full order cancellation following an earlier partial cancellation refunds `order.FinalAmount` without subtracting previous partial refunds.
  3. **Excessive Refund on Partial Returns:** Admin approving a return on a `partially_returned` order refunds `order.FinalAmount` (100% of order total) rather than the value of the returned items.

---

## 7. Security Audit (OWASP Top 10 Assessment)

### OWASP Top 10 / API Top 10 Vulnerability Matrix

| Vulnerability Category | Status | Codebase Assessment |
| :--- | :---: | :--- |
| **A01: Broken Access Control** | **FAILED** | Public IDOR on `GET /order-success/:order_id` exposes full order records. Blocked users can continue accessing protected endpoints using active JWTs. |
| **A02: Cryptographic Failures** | **FAILED** | PRNG OTP generation using `math/rand` seeded by Unix timestamp; unencrypted DB connection (`sslmode=disable`). |
| **A03: Injection** | **PASSED** | GORM parameterized queries prevent SQL injection across handlers; no string concatenation found in raw SQL queries. |
| **A04: Insecure Design** | **FAILED** | Missing coupon usage tracking; non-atomic wallet checkout; missing stock reservation during payment gateway redirect; missing current password verification on email change. |
| **A05: Security Misconfiguration** | **FAILED** | Missing CORS headers; missing security headers (HSTS, CSP, X-Frame-Options); world-writable file permissions (`0777`); Swagger documentation path mismatches. |
| **A06: Vulnerable & Outdated Components**| **FAILED** | Imports deprecated `github.com/golang-jwt/jwt` v3 alongside v5; imports deprecated `go-playground/validator` v9 alongside v10. |
| **A07: Identification & Auth Failures** | **FAILED** | Hardcoded OAuth `state="state"` allows OAuth CSRF; zero rate limiting on login, signup, and OTP endpoints; OTPs logged to console in plaintext; non-expiring email change OTPs. |
| **A08: Software & Data Integrity Failures**| **FAILED** | Webhook verification is structurally sound, but internal state transitions fail to handle out-of-stock scenarios gracefully, leaving captured payments in unhandled states. |
| **A09: Security Logging & Monitoring** | **FAILED** | Sensitive OTP credentials logged to stdout; no security event logging for failed logins, account blocks, or permission violations; no health check endpoint. |
| **A10: Server-Side Request Forgery (SSRF)**| **PASSED** | No arbitrary external URL fetching based on user input. |

---

## 8. Authentication and Authorization Audit

### JWT Implementation Details
1. **User Token Structure (`helper/helper.go`):**
   - Access Token Lifetime: 15 minutes.
   - Refresh Token Lifetime: 90 days (`24 * 90 * time.Hour`).
   - Signing Method: HMAC-SHA256 (`jwt.SigningMethodHS256`).
   - Secret Key: Loaded from `cfg.KEY`.
2. **Admin Token Structure (`helper/admin.go`):**
   - Token Lifetime: 48 hours.
   - Signing Method: HMAC-SHA256 (`jwt.SigningMethodHS256`).
   - Secret Key: Loaded from `cfg.KEY_FOR_ADMIN`.
   - **Vulnerability:** Built with deprecated `github.com/golang-jwt/jwt` v3 utilizing legacy `jwt.StandardClaims`.

### Middleware Security Gaps
1. **Token Invalidation Scope:**
   `middleware.AuthMiddleware` properly verifies that tokens are not blacklisted via `repository.IsTokenBlacklist(tokenString)`. However, tokens are only added to the blacklist when the user explicitly calls `/logout`. When an admin blocks a user (`/admin/block-user/:id`), active tokens are **not** blacklisted, allowing blocked users to continue making authenticated requests.
2. **Missing Database User Re-Verification:**
   The auth middleware decodes the JWT claims and sets `user_id` on the Gin context without verifying whether the user record still exists, has been deleted, or has been blocked in the database.
3. **Cookie Mismatch:**
   `handlers/user.go:UserLoginWithPassword` sets cookie `"Authorization"`, but `middleware/user.go:AuthMiddleware` attempts to read `"jwt_token"`. Browser cookie authentication is non-functional.
4. **Admin Role Verification:**
   `middleware.AuthorizationMiddleware` only verifies that the token was signed with `KEY_FOR_ADMIN`. It does not extract or set admin identity on `c.Set()`, preventing downstream audit logging.

---

## 9. Razorpay Payment Integration Audit

### Payment Lifecycle Verification

$$\text{Client} \xrightarrow{\text{POST /checkout/razorpay}} \text{Create Pending Order} \xrightarrow{\text{Razorpay Order API}} \text{Client Gateway} \xrightarrow{\text{Payment}} \begin{cases} \text{POST /checkout/razorpay/verify} \\ \text{POST /webhooks/razorpay} \end{cases}$$

### Detailed Assessment

1. **Order Creation (`usecase/payment.go:CreateRazorpayOrder`):**
   - **Amount Calculation:** Amount is calculated via `GetCheckoutPage` and converted to paise (`math.Round(checkout.FinalAmount * 100)`). Currency is hardcoded to `"INR"`.
   - **Flaw:** Inventory is not locked or reserved upon Razorpay order creation.
   - **Flaw:** Lacks idempotency. If a customer clicks "Pay" multiple times, multiple distinct pending orders and multiple Razorpay gateway orders are created for the same cart.
2. **Payment Verification (`usecase/payment.go:VerifyPayment`):**
   - **Signature Verification:** Reconstructs payload `payment.RazorpayOrderID + "|" + req.RazorpayPaymentID` and computes HMAC-SHA256 using `config.AppConfig.RAZORPAY_KEY_SECRET`. Compares signatures using constant-time comparison `hmac.Equal`.
   - **Ownership Check:** Correctly validates `payment.UserID == userID`.
   - **Idempotency:** Checks `if currentPayment.Status == "paid" { return nil }`.
   - **Critical Defect:** If any product or variant in the order has run out of stock while the user was completing payment, `ReduceStockAfterPaymentTx` returns an error, rolling back the database transaction. The order remains "pending", the payment remains "pending", the customer has already been charged by Razorpay, and no refund or alert is generated.
3. **Payment Failure Handler (`usecase/payment.go:HandlePaymentFailure`):**
   - Fetches the payment entity directly from Razorpay via `config.RazorpayClient.Payment.Fetch` to confirm `status == "failed"` before updating local status to failed. Prevents client-spoofed failure requests.

---

## 10. Razorpay Webhook Audit

### Webhook Verification Workflow (`handlers/payment.go:RazorpayWebhook`)
1. **Secret Validation:** Reads `config.AppConfig.RAZORPAY_WEBHOOK_SECRET`.
2. **Body Reading & Size Limit:** Reads raw body using `io.LimitReader(c.Request.Body, (1<<20)+1)` (1MB limit). Rejects payloads larger than 1MB with `413 Request Entity Too Large`.
3. **HMAC Signature Validation:**
   ```go
   mac := hmac.New(sha256.New, []byte(webhookSecret))
   mac.Write(rawBody)
   expectedSignature := mac.Sum(nil)
   receivedBytes, err := hex.DecodeString(receivedSignature)
   if err != nil || !hmac.Equal(receivedBytes, expectedSignature) {
       c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid webhook signature"})
       return
   }
   ```
   **Evaluation:** Signature validation is cryptographically correct, uses raw request bytes, and compares using constant-time `hmac.Equal`.
4. **Event Handling (`usecase/payment.go:HandleRazorpayWebhook`):**
   - Supports `payment.captured` and `payment.failed`.
   - Validates currency is `"INR"`.
   - Validates paise amount matches expected amount in database.
   - Updates payment and order status in an atomic transaction with row locking (`clause.Locking{Strength: "UPDATE"}`).
   - **Critical Defect:** If stock is unavailable, `ReduceStockAfterPaymentTx` fails, rolling back the transaction and returning an error. The handler responds with `500 Internal Server Error`, triggering Razorpay to retry the webhook repeatedly up to 24 hours while leaving the user's payment unrecorded.

---

## 11. PostgreSQL and GORM Audit

### Schema & Entity Integrity
1. **Redundant Primary Keys:**
   `domain.User` and `domain.Address` embed `gorm.Model` (which includes `ID uint gorm:"primarykey"`) and then declare `ID uint` again. While GORM resolves this, it is redundant and error-prone.
2. **Missing Database Constraints:**
   - `cart_items`: Lacks a unique composite index on `(cart_id, product_id, variant_id)`. Concurrent additions can produce duplicate rows for the same cart item.
   - `admins`: `domain.Admin` fields `Email` and `Password` use `gorm:"validate:required"`. `validate` is a struct validator tag, not a GORM schema tag. No `NOT NULL` constraint is created in PostgreSQL.
3. **Unmigrated Entity:**
   `domain.Referral` is defined in Go code but missing from `DB.AutoMigrate()`.
4. **Unindexed High-Traffic Tables:**
   `blacklist_tokens`: Checked on every authenticated request via `IsTokenBlacklist`. While `token` has a unique constraint, tokens are never purged upon expiration, leading to unbounded table growth.

### Connection Pooling & Driver Configuration
In `database/database.go`:
```go
db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Silent),
})
```
- **Missing Connection Pool Configuration:** `sqlDB, err := db.DB()` is never called to configure `SetMaxOpenConns`, `SetMaxIdleConns`, or `SetConnMaxLifetime`. In production under concurrent traffic, PostgreSQL will either exhaust available connection slots or accumulate stale connections.
- **Unencrypted Transport:** Connection string hardcodes `sslmode=disable`.

---

## 12. Go Code Quality Audit

1. **Deadlock / Startup Lockout (`main.go:57-72`):**
   Server routes are enclosed within an inverted proxy check.
2. **Ignored Errors:**
   Widespread use of `_ = ...` suppresses critical errors:
   - `usecase/cart.go:81,96`: Wishlist item deletion errors ignored.
   - `usecase/otp.go:189`: OTP deletion errors ignored.
   - `usecase/user.go:201`: Email change OTP deletion errors ignored.
   - `database/database.go:30-49`: 19 calls to `DB.AutoMigrate` ignore returned migration errors.
3. **Configuration Reloading Overhead:**
   Calling `config.LoadConfig()` inside request middleware and helper functions forces file reads and Viper YAML/env parsing on hot paths.
4. **Unchecked Type Assertions:**
   `userIDVal, _ := c.Get("user_id"); userID := userIDVal.(uint)` in `handlers/user.go` risks runtime panics if context values are corrupted or missing.
5. **Deprecated Standard Library Usage:**
   `rand.Seed(time.Now().UnixNano())` in `helper/otp.go` is deprecated since Go 1.20 and causes global mutex lock contention across goroutines.

---

## 13. Performance and Scalability Audit

1. **Unbounded Database Queries (Memory Exhaustion):**
   - `repository/sales_report.go:20`: Fetches all matching orders without pagination. Requesting a multi-month or annual sales report on a store with high transaction volume will load hundreds of thousands of full GORM models into RAM, crashing the server with Out-Of-Memory (OOM).
   - `repository/product.go:220`: If `page=0` is passed, negative offset is submitted to PostgreSQL, causing query failure.
2. **Synchronous SMTP Blocking:**
   `helper.SendOTPEmail` performs a synchronous connection to Gmail SMTP over port 587. Network latency or SMTP throttling directly blocks HTTP handlers for 2 to 5 seconds per request.
3. **N+1 Query Avoidance:**
   Preloading is appropriately utilized in `cartItems`, `products`, and `orders`. However, `GetCheckoutPage` runs multiple separate database calls sequentially rather than fetching in a single batched join.

---

## 14. Error Handling and Logging Audit

1. **Sensitive Information Leakage in Logs:**
   - `usecase/user.go:64`: `log.Println("OTP sent", otp)`
   - `usecase/otp.go:135`: `log.Println("Resend OTP:", otp)`
   OTPs must never be logged to stdout in production environments.
2. **Status Code Misuse (502 Bad Gateway):**
   Multiple handlers return `502 Bad Gateway` for client input validation failures (`ShouldBindJSON`) and password mismatches. 502 indicates that an upstream gateway/proxy server failed; client validation failures must return `400 Bad Request`, and authentication failures must return `401 Unauthorized`.
3. **Inconsistent Response Envelope:**
   `handlers/otp.go:ForgotPassword` returns raw JSON strings rather than the standardized `response.Response` struct returned by all other endpoints.

---

## 15. Health Check Endpoint Assessment and Recommendation

### Current State
**Does a health check route already exist?**
No. An exhaustive search across all files, router definitions, and handlers confirms that **no health check endpoint exists** in the ShoeGo codebase.

### Recommended Implementation Specification

#### 1. Endpoint Route & Method
- **Method:** `GET`
- **Path:** `/health`
- **Access Level:** Public (no authentication required) to allow browser access, uptime monitors, load balancers, and Kubernetes liveness/readiness probes.
- **Security Considerations:** The endpoint must reveal only status classifications (`healthy` / `unhealthy`). It must never expose connection strings, database credentials, hostnames, query latency measurements, or stack traces.

#### 2. Expected Response Payloads & Status Codes

##### Case A: System Healthy (HTTP 200 OK)
```json
{
  "status": "healthy",
  "service": "ShoeGo Backend",
  "timestamp": "2026-11-12T10:30:00Z",
  "checks": {
    "application": "healthy",
    "database": "healthy"
  }
}
```

##### Case B: Database Dependency Unhealthy (HTTP 503 Service Unavailable)
```json
{
  "status": "unhealthy",
  "service": "ShoeGo Backend",
  "timestamp": "2026-11-12T10:30:00Z",
  "checks": {
    "application": "healthy",
    "database": "unhealthy"
  }
}
```

#### 3. Verification & Testing Procedure
- **Via Web Browser:** Navigate to `http://localhost:8080/health`. A valid JSON payload with `"status": "healthy"` should display directly in the browser.
- **Via cURL:**
  ```bash
  curl -i http://localhost:8080/health
  ```
  Expected output includes `HTTP/1.1 200 OK` and the JSON body above.

#### 4. Implementation Blueprint (Files Requiring Changes When Approved)

> [!IMPORTANT]
> The following code snippets are design recommendations for future implementation. In accordance with read-only audit rules, they have NOT been applied to the codebase.

##### File 1: Create `handlers/health.go`
```go
package handlers

import (
	"context"
	"net/http"
	"time"

	"shoego/database"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func HealthCheck(c *gin.Context) {
	checks := map[string]string{
		"application": "healthy",
		"database":    "healthy",
	}

	statusCode := http.StatusOK
	overallStatus := "healthy"

	// Lightweight PostgreSQL connectivity check with 2s timeout
	sqlDB, err := database.DB.DB()
	if err != nil {
		checks["database"] = "unhealthy"
		overallStatus = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			checks["database"] = "unhealthy"
			overallStatus = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}
	}

	c.JSON(statusCode, HealthCheckResponse{
		Status:    overallStatus,
		Service:   "ShoeGo Backend",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
```

##### File 2: Register in `router/user.go` (or `main.go`)
```go
// Register public health check route
r.GET("/health", handlers.HealthCheck)
```

---

## 16. Testing and Verification Results

### Test Execution Summary
- **Command:** `go test ./...`
- **Result:**
  ```
  ?   	shoego	[no test files]
  ?   	shoego/config	[no test files]
  ?   	shoego/database	[no test files]
  ?   	shoego/docs	[no test files]
  ?   	shoego/domain	[no test files]
  ?   	shoego/handlers	[no test files]
  ?   	shoego/helper	[no test files]
  ?   	shoego/middleware	[no test files]
  ?   	shoego/models	[no test files]
  ?   	shoego/repository	[no test files]
  ?   	shoego/response	[no test files]
  ?   	shoego/router	[no test files]
  ?   	shoego/usecase	[no test files]
  ```
- **Tests Executed:** `0`
- **Tests Passed:** `0`
- **Tests Failed:** `0`
- **Unit Test Coverage:** `0.0%`

### Static Code Analysis Results
- **Command:** `go vet ./...`
- **Result:** Exited with code 0 (no syntax/vet violations reported).
- **Compilation Check:** `go build -o /dev/null main.go` compiles successfully, confirming that the startup bug in `main.go` is a semantic control flow defect rather than a compiler error.

---

## 17. Critical Issues

### [CRIT-01] Fatal Server Startup Flow Control Defect in `main.go`
- **Issue ID:** CRIT-01
- **Severity:** Critical
- **Category:** Architecture / Logic
- **File Location:** `main.go`
- **Line Number:** Lines 57–72
- **Description:** Route registrations, static folder serving, and `router.Run("0.0.0.0:8080")` are nested inside `if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil { ... }`.
- **Impact:** When `SetTrustedProxies` succeeds, `main()` exits immediately without listening on port 8080. If it fails, `log.Fatal` terminates the process. The server cannot run in production.
- **Verification Status:** Verified via static code inspection.
- **Reproduction:** Inspect lines 57-72 of `main.go`. Notice that line 57 opens an error check block, and line 72 closes it after `router.Run`.
- **Recommendation:** Move all route registration, static file serving, swagger setup, and `router.Run` outside the `if err != nil` block.

---

### [CRIT-02] Public Unauthenticated IDOR Order Details Exposure
- **Issue ID:** CRIT-02
- **Severity:** Critical
- **Category:** Security / Broken Access Control
- **File Location:** `handlers/checkout.go` & `router/user.go`
- **Line Number:** `handlers/checkout.go:235-252`, `router/user.go:27`
- **Description:** `GET /order-success/:order_id` is registered as a public unauthenticated route. Inside `handlers/checkout.go`, the ownership check `userID, ok := getUserIDFromContext(c)` is commented out. `usecase.GetOrderSuccess` calls `repository.GetOrderByOrderIDPayment(orderID)` without filtering by user ID.
- **Impact:** Anyone can access the full order confirmation data (order items, delivery address, phone number, customer name, payment amount) for any valid order ID.
- **Verification Status:** Verified via static code inspection.
- **Reproduction:** Call `GET /order-success/<any_order_id>` without passing an `Authorization` header.
- **Recommendation:** Move `/order-success/:order_id` under `userProtected` middleware group and verify that `order.UserID == authenticatedUserID`.

---

### [CRIT-03] Unlimited Coupon Usage Due to Missing Usage Increment
- **Issue ID:** CRIT-03
- **Severity:** Critical
- **Category:** Logic / Financial
- **File Location:** `models/coupon.go:20`, `usecase/coupon.go:122,237`, `domain/coupon.go:17`
- **Line Number:** Entire codebase
- **Description:** While `coupon.UsedCount >= coupon.UsageLimit` is checked during coupon application, `coupon.UsedCount` is never incremented anywhere in the codebase when an order is finalized.
- **Impact:** Limited-use coupons (e.g. valid for first 50 users or single-use promo codes) can be used infinitely by unlimited users, causing financial loss.
- **Verification Status:** Verified via codebase grep (`UsedCount` is never updated or written to database).
- **Reproduction:** Apply a coupon with `UsageLimit: 1`, place an order, and observe that `coupons.used_count` remains 0. Apply the same coupon to a second order; it succeeds.
- **Recommendation:** In `CreateOrderTransaction` and Razorpay payment success transactions, execute an atomic SQL update: `UPDATE coupons SET used_count = used_count + 1 WHERE code = ? AND used_count < usage_limit`. Create a `user_coupons` table to track per-user usage.

---

### [CRIT-04] Wallet Double-Spending via Missing Locks & Non-Atomic Debit
- **Issue ID:** CRIT-04
- **Severity:** Critical
- **Category:** Database / Financial / Concurrency
- **File Location:** `usecase/checkout.go`
- **Line Number:** Lines 302–386
- **Description:** `PlaceWalletOrder` checks `wallet.Balance < checkout.FinalAmount` on a record fetched without row locking (`FOR UPDATE`). In-memory balance calculation (`newBalance := wallet.Balance - checkout.FinalAmount`) is written back via a separate update query outside the order placement transaction.
- **Impact:** Two concurrent checkout requests can spend the same wallet balance twice, allowing users to obtain free products.
- **Verification Status:** Verified via static code inspection.
- **Reproduction:** Trigger two concurrent POST requests to `/checkout/wallet` with a wallet balance equal to one order's value. Both succeed.
- **Recommendation:** Enclose order creation and wallet deduction in a single database transaction. Use row locking (`clause.Locking{Strength: "UPDATE"}`) and atomic SQL deduction: `UPDATE wallets SET balance = balance - ? WHERE id = ? AND balance >= ?`.

---

### [CRIT-05] Missing Payment Record in `PlaceWalletOrder` Forfeits Refund on Cancellation
- **Issue ID:** CRIT-05
- **Severity:** Critical
- **Category:** Logic / Financial
- **File Location:** `usecase/checkout.go` & `usecase/order.go`
- **Line Number:** `usecase/checkout.go:350-385`, `usecase/order.go:348-359`
- **Description:** `PlaceWalletOrder` creates an order with `PaymentMethod: "wallet"`, but fails to insert a record into the `payments` table. When the customer cancels the order, `CancelOrder` checks `payment, err := repository.GetPaymentByOrderID(order.ID)`. Since no payment row exists, `GetPaymentByOrderID` returns an error, causing the refund check `if err == nil && payment.Status == "paid"` to evaluate to false.
- **Impact:** Customers who pay with their wallet and subsequently cancel their order lose 100% of their money with zero wallet refund.
- **Verification Status:** Verified via code path tracing.
- **Reproduction:** Place an order via `/checkout/wallet`, then call `/orders/:order_id/cancel`. Observe that wallet balance does not increase.
- **Recommendation:** Insert a `domain.Payment` record with `Status: "paid"`, `PaymentMethod: "wallet"` inside `CreateOrderTransaction`, or update `CancelOrder` to inspect `order.PaymentMethod == "wallet"` directly.

---

### [CRIT-06] Stock Depletion During Razorpay Payment Causes Stuck Funds and 500 Webhook Loop
- **Issue ID:** CRIT-06
- **Severity:** Critical
- **Category:** Payment / Razorpay
- **File Location:** `usecase/payment.go` & `repository/payment.go`
- **Line Number:** `usecase/payment.go:222-229`, `repository/payment.go:216-248`
- **Description:** `CreateRazorpayOrder` creates an order without reserving inventory. If stock is depleted by another customer before payment completes, `ReduceStockAfterPaymentTx` returns `variant stock unavailable`. This rolls back `ProcessWebhookPaymentSuccessTx`.
- **Impact:** Razorpay captures the customer's payment, but ShoeGo rolls back the database update. The order remains "pending", the user receives no goods, no automated refund is triggered, and the webhook returns 500, causing Razorpay to retry endlessly.
- **Verification Status:** Verified via transaction rollback inspection.
- **Recommendation:** If stock reduction fails during payment verification/webhook, catch the error, mark the order as `stock_exhausted`, record the payment as `paid`, and immediately invoke Razorpay Refund API or credit the user's wallet with an automated refund.

---

### [CRIT-07] Hardcoded OAuth State Parameter and Missing CSRF Validation
- **Issue ID:** CRIT-07
- **Severity:** Critical
- **Category:** Security / Authentication
- **File Location:** `handlers/logingoogle.go`
- **Line Number:** Lines 25, 43–47
- **Description:** `GoogleLogin` generates authentication URLs using a static hardcoded state: `config.GoogleOAuthConfig.AuthCodeURL("state")`. In `GoogleCallback`, the `state` parameter is never checked against a user session or cookie.
- **Impact:** Vulnerable to OAuth Cross-Site Request Forgery (CSRF). An attacker can trick a victim into linking the attacker's Google account or logging into the attacker's session.
- **Verification Status:** Verified via code inspection.
- **Recommendation:** Generate a cryptographically random state token (`crypto/rand`), store it in a secure HTTP-only cookie, and verify it upon callback before exchanging the authorization code.

---

### [CRIT-08] Blocked User Access via Google OAuth and Active JWTs
- **Issue ID:** CRIT-08
- **Severity:** Critical
- **Category:** Security / Authorization
- **File Location:** `handlers/logingoogle.go` & `middleware/user.go`
- **Line Number:** `handlers/logingoogle.go:55-83`, `middleware/user.go:38-85`
- **Description:** When an account is blocked (`blocked = true`), `GoogleCallback` does not verify `existingUser.Blocked` before issuing tokens. Furthermore, `middleware.AuthMiddleware` only checks if the token string is in `blacklist_tokens` and never queries the database to verify if the account has been blocked.
- **Impact:** Administrative blocking is completely ineffective. Blocked users can continue to use active JWTs or simply log in via Google OAuth.
- **Verification Status:** Verified via code inspection.
- **Recommendation:** In `GoogleCallback`, check `if existingUser.Blocked { return Unauthorized }`. In `AuthMiddleware`, verify user account status or blacklist all active tokens for that user ID when blocked.

---

## 18. High-Severity Issues

- **[HIGH-01] Plaintext OTP Logged to Server Console (`usecase/user.go:64`, `usecase/otp.go:135`):** `log.Println("OTP sent", otp)` prints raw OTPs to stdout. Anyone with access to system logs or log aggregation services can intercept OTPs and hijack accounts.
- **[HIGH-02] Insecure PRNG Used for Security Tokens (`helper/otp.go:9-16`):** Uses `math/rand` with timestamp seed instead of `crypto/rand`. Predictable OTP generation allows account takeover.
- **[HIGH-03] Order Cancellation Corrupts Stock Duality (`repository/order.go:317-339`):** When cancelling a variant order item, it increments both parent `Product.Stock` AND `ProductVariant.Stock`, generating phantom inventory.
- **[HIGH-04] Double Refund on Multi-Step Cancellations (`usecase/order.go:346-359`):** `CancelOrder` refunds `order.FinalAmount` without checking whether individual order items were already cancelled and refunded earlier.
- **[HIGH-05] Non-Expiring Email Change OTPs (`usecase/user.go:173-203`):** `RequestEmailChange` creates OTP records without setting `ExpiresAt`, and `VerifyEmailChange` never validates expiration.
- **[HIGH-06] Full Refund Granted on Partial Order Returns (`usecase/order.go:148-164`):** When admin updates an order status from `partially_returned` to `returned`, `CreditWallet` credits `order.FinalAmount` (full order value).
- **[HIGH-07] Cookie Name Mismatch Between Login and Middleware (`handlers/user.go:100`, `middleware/user.go:19`):** Login sets cookie `"Authorization"`, middleware reads `"jwt_token"`. Cookie auth fails entirely.
- **[HIGH-08] Zero Rate Limiting on Authentication Endpoints (`router/user.go`, `router/admin.go`):** Login, signup, password reset, and OTP verification endpoints are unprotected against brute-force attacks.
- **[HIGH-09] Cart Operations Affect Multiple Variants of Same Product (`repository/cart.go:95, 114-123`):** Updating or removing an item using `product_id` deletes or modifies all sizes/variants of that product in the cart.
- **[HIGH-10] Permanent Sticky Coupons on Non-Razorpay Orders (`repository/checkout.go:88-100`):** `CreateOrderTransaction` does not clear `cart.CouponCode` on COD and Wallet orders, applying the coupon to future purchases automatically.
- **[HIGH-11] Unbounded File Uploads & Insecure Permissions (`helper/products.go:30`):** Directory created with `0777` permissions; lack of maximum upload size check enables denial of service via memory exhaustion.
- **[HIGH-12] Missing Database Connection Pool Configuration (`database/database.go:16-53`):** No limits on open/idle connections; 19 calls to `DB.AutoMigrate` ignore returned errors.
- **[HIGH-13] Unmigrated Domain Model (`domain/offer.go:32-38`, `database/database.go:30-49`):** `domain.Referral` is defined but omitted from `DB.AutoMigrate`.

---

## 19. Medium-Severity Issues

- **[MED-01] Config Loaded from Disk on Every Request (`helper/auth.go:21`, `helper/admin.go:20`):** Re-reading `.env` inside token verification degrades API throughput.
- **[MED-02] Missing Viper Env Bindings for Secrets (`config/config.go:39-40`):** `KEY` and `KEY_FOR_ADMIN` are missing from `viper.BindEnv`, breaking containerized deployments that pass secrets via OS environment variables without a `.env` file.
- **[MED-03] Hardcoded `sslmode=disable` (`database/database.go:17`):** Database connection traffic is unencrypted across the network.
- **[MED-04] Misuse of HTTP 502 Bad Gateway (`handlers/admin.go:28-51`, `handlers/user.go:71-92`):** Client validation errors and password mismatches return 502 instead of 400 or 401.
- **[MED-05] Inconsistent Non-Envelope Response in `ForgotPassword` (`handlers/otp.go:94,98`):** Returns raw string instead of `response.ClientResponse`.
- **[MED-06] In-Memory Map Used for Checkout Handoff (`handlers/checkout.go:25-28`):** Handoff links break when running behind multi-instance load balancers.
- **[MED-07] Unchecked Type Assertions (`handlers/user.go:144,166,188`):** `c.Get("user_id").(uint)` can cause panics if the key is missing or invalid.
- **[MED-08] Unbounded In-Memory Query in Sales Report (`repository/sales_report.go:20`):** Loads all matching orders into memory without pagination.
- **[MED-09] Negative Offset Query Crash (`repository/product.go:220`):** If `page=0` is passed, negative offset is submitted to PostgreSQL.
- **[MED-10] Non-Variant Stock Reports 0 in Cart (`usecase/cart.go:206`):** `item.Variant.Stock` returns 0 for non-variant products in `GetCart`.
- **[MED-11] Unindexed Blacklist Table Growth (`domain/user.go:31-34`):** Blacklisted tokens are never purged upon expiration.
- **[MED-12] Dual JWT Dependencies (`go.mod:10,11`):** Imports deprecated `golang-jwt/jwt` v3 alongside v5.

---

## 20. Low-Severity Issues

- **[LOW-01] Suppressed Error Returns (`usecase/cart.go:81`, `usecase/otp.go:189`):** Discarding errors with `_ =` masks silent failures.
- **[LOW-02] Missing CORS Configuration (`main.go`):** Browsers block single-page applications hosted on distinct origins.
- **[LOW-03] Missing Standard Security Headers (`main.go`):** Lacks HSTS, CSP, X-Frame-Options, and X-Content-Type-Options.
- **[LOW-04] Inconsistent Query Parameters (`handlers/sales_reports.go`):** `filter, start_date, end_date` vs `period, from, to`.
- **[LOW-05] Asymmetric REST Route Naming (`router/admin.go:48,49`):** PUT `/offers/:id` vs DELETE `/offers/product/:id`.
- **[LOW-06] Typographical Errors in Responses:** Errors such as `"formate is not correct admin"` and `"constrian are not satisfied"`.
- **[LOW-07] Redundant Primary Key Declarations (`domain/user.go:6-7`):** Re-declaring `ID uint` after embedding `gorm.Model`.

---

## 21. Informational Improvements

- **[INFO-01] Zero Test Coverage:** Codebase contains zero unit, integration, or end-to-end tests.
- **[INFO-02] Synchronous SMTP Operations:** Offload email delivery to background worker goroutines or asynchronous message queues (RabbitMQ/Redis).
- **[INFO-03] Missing Structured Logging:** Replace standard `log` and `fmt.Println` with `uber-go/zap` or Go's standard `log/slog`.
- **[INFO-04] Swagger Documentation Drift:** Documented paths (e.g. `/admin/admin`) do not match registered routes (`/admin/login`).

---

## 22. Recommended Fix Priority & Remediation Roadmap

```
Phase 1: Fatal Startup & Security Kill-Chain
  ├── Fix main.go inverted proxy check block (CRIT-01)
  ├── Protect /order-success/:order_id IDOR with auth middleware & user check (CRIT-02)
  ├── Implement secure cryptographic OAuth state validation (CRIT-07)
  ├── Enforce account block checks in OAuth and AuthMiddleware (CRIT-08)
  └── Remove plaintext OTP logging from stdout (HIGH-01)

Phase 2: Financial & Transaction Integrity
  ├── Implement atomic coupon usage incrementing & user-coupon tracking (CRIT-03)
  ├── Wrap PlaceWalletOrder in single transaction with FOR UPDATE row lock (CRIT-04)
  ├── Insert Payment record for wallet checkouts to fix refund bug (CRIT-05)
  ├── Implement stock reservation or auto-refund on Razorpay out-of-stock (CRIT-06)
  └── Correct stock restoration logic on item/order cancellations (HIGH-03, HIGH-04)

Phase 3: Functional & Concurrency Fixes
  ├── Replace math/rand with crypto/rand for OTP generation (HIGH-02)
  ├── Standardize cookie names between handlers and middleware (HIGH-07)
  ├── Scope cart item updates/deletions by item ID instead of product ID (HIGH-09)
  ├── Clear cart.CouponCode on COD and Wallet order placement (HIGH-10)
  └── Implement Redis or distributed store for checkout handoffs (MED-06)

Phase 4: Database & Infrastructure Hardening
  ├── Configure PostgreSQL connection pooling (MaxOpen, MaxIdle, Lifetime) (HIGH-12)
  ├── Enable SSL mode for PostgreSQL connection (MED-03)
  ├── Add Rate Limiting middleware on authentication routes (HIGH-08)
  ├── Bind all environment variables in Viper config loader (MED-02)
  └── Migrate domain.Referral and clean up domain.Admin (HIGH-13)

Phase 5: Code Quality & Observability
  ├── Implement GET /health check endpoint (Section 15)
  ├── Establish unit test suite achieving >80% coverage on usecases
  ├── Eliminate duplicate dependencies (remove jwt v3, validator v9) (MED-12)
  └── Adopt structured logging (slog/zap) and standardize API error codes (MED-04)
```

---

## 23. Final Assessment

### Quantitative Audit Summary

| Metric | Measured Value |
| :--- | :---: |
| **Total Identified Issues** | **44** |
| Critical Issues | 8 |
| High-Severity Issues | 13 |
| Medium-Severity Issues | 12 |
| Low-Severity Issues | 7 |
| Informational Improvements | 4 |
| **Total Routes Inspected** | **96** (93 API routes + 3 static/doc endpoints) |
| **Features Inspected** | **25+** (Auth, Products, Variants, Carts, Wishlists, Checkout, Orders, Returns, Payments, Razorpay, Coupons, Offers, Wallet, Reports, Dashboard) |
| **Total Tests Executed** | **0** |
| Tests Passed | 0 |
| Tests Failed | 0 |
| **Areas That Could Not Be Verified at Runtime** | Live Razorpay gateway callbacks, live Google OAuth token exchanges, live Gmail SMTP delivery (no real credentials/funds executed in read-only audit). |

---

### Non-Modification Confirmation

In strict compliance with the audit rules:
> **Explicit Confirmation:** No existing codebase files, schemas, configurations, business logic, routes, or dependencies were altered, created, deleted, or refactored during this audit. The working tree remains completely clean. The only new file created is this audit report: `CODEBASE_AUDIT_REPORT.md`.
