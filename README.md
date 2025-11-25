# Mini-Quicko

**Mini-Quicko** is a Kaspi.kz marketplace price analysis service that helps detect price dumping and optimize pricing strategies for sellers.

## Features

- 🔍 **Price Analysis** - Analyzes product prices across all merchants on Kaspi.kz
- 📊 **Statistical Insights** - Calculates min/max/average/median/optimal prices
- ⚠️ **Dumping Detection** - Identifies merchants selling 20% below average price
- 📈 **Price History** - Tracks price changes over time
- 🏪 **Seller Information** - Lists all sellers with ratings and reviews
- ⚡ **Smart Caching** - Reduces API calls with 30-minute cache validity
- 🐳 **Docker Support** - Fully containerized with Docker Compose

## Tech Stack

- **Backend**: Go 1.25.4
- **Web Framework**: Gin
- **Database**: PostgreSQL 16
- **Logging**: slog (structured JSON logging)
- **Configuration**: YAML + .env
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   │   ├── config.go
│   │   └── config.yaml
│   ├── models/                 # Data models
│   │   └── entity.go
│   ├── storage/                # Database layer (Repository pattern)
│   │   ├── storage.go          # Repository interfaces
│   │   ├── postgres.go         # Database connection
│   │   └── postgres_storage.go # PostgreSQL implementation
│   ├── service/                # Business logic layer
│   │   ├── service.go          # Service interfaces
│   │   ├── kaspi_client.go     # Kaspi.kz API client
│   │   ├── analysis_service.go # Price analysis logic
│   │   ├── history_service.go  # Price history retrieval
│   │   └── offer_fetcher.go    # Cache-first offer fetching
│   └── handlers/               # HTTP layer
│       ├── handler.go
│       ├── routes.go           # API endpoints
│       └── server.go           # Gin server setup
├── init.sql                    # Database schema
├── .env                        # Environment variables (DB password)
├── docker-compose.yml          # Docker services configuration
├── Dockerfile                  # Application container
└── Makefile                    # Build commands
```

## Architecture

The project follows a clean **Repository → Service → Handler** architecture pattern:

1. **Repository Layer** - Database operations with caching
2. **Service Layer** - Business logic and external API integration
3. **Handler Layer** - HTTP request handling with Gin

### Caching Strategy

- Offers are cached for **30 minutes** to minimize API calls
- Cache-first approach: checks database before calling Kaspi API
- Automatically saves fresh data to cache on API calls

## Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)
- PostgreSQL 16 (handled by Docker)

## Installation & Setup

### 1. Clone the repository
```bash
cd kaspi
```

### 2. Configure environment variables
Create a `.env` file:
```env
DB_PASSWORD=your_secure_password
```

### 3. Review configuration
Edit `internal/config/config.yaml` if needed:
```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  host: postgres  # Use 'localhost' for local development
  port: "5432"
  user: admin
  dbname: mini_quicko
  sslmode: disable

kaspi:
  default_city_id: "710000000"  # Almaty


app:
  environment: development
  log_level: info  # debug, info, warn, error
```

### 4. Build and run
```bash
# Build Docker images
make build

# Start all services
make up

# Stop services
make stop

# Run database migrations
make migrate
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-25T12:00:00Z",
  "version": "1.0.0"
}
```

---

### Analyze Product Prices
```http
POST /api/v1/analyze
Content-Type: application/json

{
  "product_id": "102298404",
  "city_id": "710000000"  // Optional, defaults to config
}
```
**Response:**
```json
{
  "product_id": "102298404",
  "min_price": 150000.00,
  "max_price": 200000.00,
  "avg_price": 175000.00,
  "median_price": 170000.00,
  "optimal_price": 178500.00,
  "total_offers": 25,
  "dumping_merchants": [
    {
      "merchant_id": "12345",
      "merchant_name": "TechStore KZ",
      "price": 130000.00,
      "below_avg_percent": 25.71,
      "rating": 4.8,
      "reviews_count": 1234
    }
  ],
  "timestamp": "2025-11-25T12:00:00Z"
}
```

---

### Get All Sellers
```http
GET /api/v1/products/{product_id}/sellers?city_id=710000000
```
**Response:**
```json
{
  "product_id": "102298404",
  "sellers": [
    {
      "merchant_id": "12345",
      "merchant_name": "TechStore KZ",
      "rating": 4.8,
      "reviews_count": 1234,
      "price": 175000.00
    }
  ],
  "count": 25
}
```

---

### Get Dumping Sellers
```http
GET /api/v1/products/{product_id}/dumping-sellers?city_id=710000000
```
**Response:**
```json
{
  "product_id": "102298404",
  "average_price": 175000.00,
  "dumping_sellers": [
    {
      "merchant_id": "12345",
      "merchant_name": "TechStore KZ",
      "rating": 4.8,
      "reviews_count": 1234,
      "price": 130000.00
    }
  ],
  "dumping_count": 3
}
```

---

### Get Price History
```http
GET /api/v1/history/{product_id}
```
**Response:**
```json
{
  "product_id": "102298404",
  "history": [
    {
      "product_id": "102298404",
      "product_name": "iPhone 13 128GB",
      "timestamp": "2025-11-25T12:00:00Z",
      "min_price": 150000.00,
      "max_price": 200000.00,
      "avg_price": 175000.00,
      "offer_count": 25
    }
  ]
}
```

---

### Get Raw Offers
```http
GET /api/v1/products/{product_id}/offers?city_id=710000000
```
**Response:**
```json
{
  "offers": [...],
  "total": 50,
  "offersCount": 25,
  "badges": [...],
  "highRatingPresent": true,
  "excellentMerchantPresent": true
}
```

## Finding Product IDs

1. Go to [kaspi.kz](https://kaspi.kz)
2. Search for any product
3. Copy the number from URL: `kaspi.kz/shop/p/product-name-**102298404**`
4. Use this ID in API requests

## Database Schema

### `price_history` Table
Stores historical price data for products:
- `product_id` - Kaspi product identifier
- `product_name` - Product title
- `min_price`, `max_price`, `avg_price` - Price statistics
- `offer_count` - Number of offers at that time
- `timestamp` - When data was recorded

### `offers_cache` Table
Caches offer data to reduce API calls:
- `product_id`, `city_id` - Cache key
- Merchant information (ID, name, rating, reviews)
- `price` - Merchant's price
- `fetched_at` - Cache timestamp (30-minute validity)

## Development

### Run locally without Docker
```bash
# Ensure PostgreSQL is running locally
# Update config.yaml: host: localhost, port: "5434"

go run cmd/main.go
```

### Database Operations
```bash
# Connect to database
docker exec -it mini-quicko-db psql -U admin -d mini_quicko

# List tables
\dt

# View table structure
\d price_history
\d offers_cache

# Query data
SELECT * FROM offers_cache WHERE product_id = '102298404';

# Exit
\q
```

### View Logs
```bash
# Application logs
docker logs -f mini-quicko-app

# Database logs
docker logs -f mini-quicko-db
```

## Configuration

### Log Levels
Set in `config.yaml`:
- `debug` - Detailed logs including SQL queries
- `info` - General application flow (default)
- `warn` - Warning messages only
- `error` - Errors only

### Cache Duration
Modify `cacheMaxAgeMinutes` in `internal/service/offer_fetcher.go` (default: 30 minutes)

### Kaspi API Settings
- `default_city_id`: City for price queries (710000000 = Almaty)


## Troubleshooting

### Port Already in Use
```bash
# Change ports in docker-compose.yml
ports:
  - "8081:8080"  # Application
  - "5435:5432"  # PostgreSQL
```

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep mini-quicko-db

# Restart services
make stop
make up
```

### Cache Not Working
```bash
# Clear cache manually
docker exec -it mini-quicko-db psql -U admin -d mini_quicko \
  -c "DELETE FROM offers_cache WHERE product_id = 'YOUR_PRODUCT_ID';"
```





## Support

For issues or questions, please check the Kaspi.kz API behavior and ensure your product IDs are valid.
