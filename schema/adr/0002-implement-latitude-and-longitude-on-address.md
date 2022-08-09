# 2. Implement latitude and longitude on address

Date: 2021-03-29

## Status

Accepted

## Context
```
The issue motivating this decision, and any context that influences or constrains the decision.
```

### Background
Saat ini data address atau alamat tidak menyimpan informasi latitude dan longitude. Informasi latitude dan longitude dibutuhkan untuk mengetahui titik koordinat lokasi alamat secara spesifik. Masalah yang sering dihadapi ketika informasi latitude dan longitude tidak tersedia, salah satunya adalah kurir Armada Bhinneka yang ketika mengantarkan paket ke customer selalu menanyakan ke customer terlebih dahulu melalui aplikasi lain (misalnya, WhatsApp) untuk melakukan "Share Location" agar barang dapat dikirimkan, hal ini berisiko jika customer sedang tidak online, dimana barang berpotensi untuk tidak terkirim atau terlambat dikirim. Selain itu behavior ini (meminta customer untuk Share Location via Whatsapp) juga tidak lazim dilakukan oleh kurir lain.

### Business Impact
Selain untuk solving case kurir Armada, informasi latitude dan longitude juga dapat digunakan untuk berbagai macam business usecase seperti integrasi dengan kurir yang mengharuskan adanya latitude dan longitude seperti Gojek, Grab serta layanan 3PL lainnya. Selain itu informasi latitude dan longitude juga dapat digunakan oleh Business Intelligence (BI) untuk mengetahui demografi customer baik buyer maupun seller.

## Decision

```
The change that we're proposing or have agreed to implement.
```

Untuk mendapatkan informasi latitude dan longitude akan dilakukan integrasi dengan menggunakan Google Maps API.

### Flowchart
![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbmJbV2ViLyBNb2JpbGUgQXBwcyBDbGllbnRdLS0-fEdldCBMYXQgTG9uZ3xhW0dvb2dsZSBNYXBzIEFQSV07XG5iW1dlYi8gTW9iaWxlIEFwcHMgQ2xpZW50XS0tPnxHcmFwaFFMfGNbR1dTIEFQSV07XG5jW0dXUyBBUEldLS0-fFJFU1R8ZFtCRSBTdHVyZ2Vvbl07XG5kW0JFIFN0dXJnZW9uXS0tPnxFbmNyeXB0ICYgc3RvcmV8ZVsoRGF0YWJhc2UpXTsiLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9LCJ1cGRhdGVFZGl0b3IiOmZhbHNlfQ)

Flow ke Google Maps API terdiri dari 2 skenario:
1. Geocoding, untuk mendapatkan latitude & longitude berdasarkan lokasi yang dicari.
2. Reverse Geocoding, kebalikan dari Geocoding, yaitu untuk mendapatkan lokasi berdasarkan latitude & longitude.  

### Database Schema
Ada dua pendekatan atau opsi ketika mendesain database schema, masing-masing opsi memiliki pros & cons.

#### Option 1 (Menggunakan existing table, tambah fields baru):
##### Affected Table
1. Table `b2c_shipping_address` -->  shipping address untuk personal account
2. Table `b2c_merchant` --> seller address (sebagai asal pengiriman)
3. Table `address` --> generic table untuk address
4. Table `b2b_address` --> address untuk corporate account (alamat perusahaan)
5. Table `b2b_contact_address` --> address untuk contact pada corporate account

##### Tambahan fields pada setiap table
| field         | type               | length |
|---------------|--------------------|--------|
| label         | string             | 50     |
| latitude      | decimal            | (8,6)  |
| longitude     | decimal            | (9,6)  |

**Pros:** 
- Simple & straightforward
- Performance

**Cons:**
- Sulit untuk di-maintain, jika kita ingin menambahkan field baru, setiap table harus di-altered schema-nya.
- Sulit untuk scale

#### Option 2 (Menggunakan table terpisah):
Buat table baru dengan nama `maps`, dengan menggunakan teknik [Polymorphism](https://docs-api.bhinneka.com/backend-tech-talk.html#polymorphism-in-go-and-postgresql).

| field         | type                              | length |
|---------------|-----------------------------------|--------|
| id (PK)       | string (e.g. MAPS201202141056740) | 50     |
| relationId    | string                            | 50     |
| relationName  | string                            | 50     |
| label         | string                            | 50     |
| latitude      | decimal                           | (8,6)  |
| longitude     | decimal                           | (9,6)  |

**Pros:** 
- Mudah untuk di-maintain, dengan menghindari redundancy
- Lebih mudah untuk scale dengan table yang generic

**Cons:**
- Performance, tapi jika semua table yang terkait dengan address sudah di-merged menjadi satu table, harusnya ini tidak menjadi isu lagi
- Tidak ada foreign key

**Notes:** All data type e.g. `decimal` will be changed to `text` if needs to be encrypted.

#### Solution:
Kita pilih opsi nomor 2 karena lebih mudah untuk di-maintain & lebih cocok untuk jangka panjang, misal nanti ingin dibuatkan POI service sendiri untuk mem-proxy API Google Maps atau koneksi ke provider map lain, maka tinggal meng-extract satu table saja yaitu `maps`.

### Affected Endpoint
1. POST /api/v2/shipping-address/me
2. PUT /api/v2/shipping-address/me/{ADDRID}
3. GET /api/v2/shipping-address/me/{ADDRID}
4. POST /api/v2/merchant/me/warehouse
5. PUT /api/v2/merchant/me/warehouse/{ADDRID}
6. PUT /api/v2/merchant/{MerchantID}
7. GET /api/v2/merchant/list (CMS Merchant)
8. GET /api/v2/merchant/{MERCHANTID} (CMS Get Merchant)
9. POST /api/v2/merchant (Merchant Data Register)
10. POST /api/v2/shipping-address (CMS Create Address)
11. PUT /api/v2/shipping-address/{ADDRID}
12. GET /api/v2/shipping-address (CMS List)
13. GET /api/v2/shipping-address/{ADDRID} (CMS Get Address)

### Sample Payload

#### 1. GET /addresses
```json
{
    "success": true,
    "code": 200,
    "message": "List Address",
    "data": [
        {
            "label": "rumah",
            "name": "john doe",
            "address": "ini adalah alamat baris 1",
            "provinceId": "0104",
            "provinceName": "Jawa Barat",
            "cityId": "010404",
            "cityName": "Bekasi",
            "districtId": "01040405",
            "districtName": "Bekasi Timur",
            "subDistrictId": "0104040501",
            "subDistrictName": "Aren Jaya",
            "postalCode": "17111",
            // other fields
            "isMapAvailable": true // boolean, to indicate latitude & longitude available (already filled)
        }
    ],
    "meta": {
        "page": 1,
        "limit": 1,
        "totalRecords": 16,
        "totalPages": 16
    }
}
```

#### 2. POST /addresses/{ADDRID}/maps
```json
{
    "map": {
        "label": "Stasiun Pasar Minggu",
        "latitude": 37.4224764,
        "longitude": -122.0842499
        // any other field
    }
}
```

#### 3. GET /addresses/{ADDRID}
```json
{
    "success": true,
    "code": 200,
    "message": "Detail Address",
    "data": {
        "label": "rumah",
        "name": "john doe",
        "address": "ini adalah alamat baris 1",
        "provinceId": "0104",
        "provinceName": "Jawa Barat",
        "cityId": "010404",
        "cityName": "Bekasi",
        "districtId": "01040405",
        "districtName": "Bekasi Timur",
        "subDistrictId": "0104040501",
        "subDistrictName": "Aren Jaya",
        "postalCode": "17111",
        // other fields
        "isMapAvailable": true, // boolean, to indicate latitude & longitude available (already filled)
        "map": {
            "id": "MAPS201202141056740",
            "label": "Stasiun Pasar Minggu",
            "latitude": 37.4224764,
            "longitude": -122.0842499
            // any other field
        }
    }
}
```
**Notes:** Data type for latitude & longitude is `double`

## Consequences

```
What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.
```

Informasi latitude dan longitude merupakan data PII (Personal Identifiable Information) sehingga diwajibkan untuk dilakukan enkripsi di database.

Selain itu dari flowchart terlihat bahwa client directly connect ke Google Maps API, sehingga hit/ traffic ke Google Maps API akan tinggi & meningkatkan cost, ke depannya dibutuhan POI service untuk mem-proxy/ caching agar direct hit dapat diminimalisir.