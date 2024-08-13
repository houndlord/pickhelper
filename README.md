# [Pickhelper](https://pickhelper.lol/)

A service which scrapes statistics data from op.gg and returns them via web frontend and REST API.

```
[op.gg] <-- Scrape -- [Go Scraper/API] <-- HTTP --> [REACT] <-- HTTP --> [User's Browser]
                            ^
                            |
                            v
                        [Database]
```

## Endpoints

To directly call endpoints: https://pickhelper.lol/api

### 1. Get All Champions

Retrieves a list of all champions.

- **URL:** `/champions`
- **Method:** `GET`
- **Success Response:**
  - **Code:** 200
  - **Content:** 
    ```json
    {
      "champions": [
        {
          "Name": "Ahri",
          "AvatarURL": "https://example.com/ahri.png"
        },
        {
          "Name": "Zed",
          "AvatarURL": "https://example.com/zed.png"
        },
        ...
      ]
    }
    ```

### 2. Get Matchups for a Champion

Retrieves matchup data for a specific champion in a specific role.

- **URL:** `/matchups/:champion/:role`
- **Method:** `GET`
- **URL Parameters:**
  - `champion`: The name of the champion
  - `role`: The role (top, jungle, mid, adc, support)
- **Query Parameters:**
  - `limit` (optional): Number of matchups to return (default: 8)
- **Success Response:**
  - **Code:** 200
  - **Content:** 
    ```json
    {
      "patch": "11.10",
      "matchups": [
        {
          "Champion": "Zed",
          "WinRate": "55.5",
          "SampleSize": "1000"
        },
        {
          "Champion": "Yasuo",
          "WinRate": "52.3",
          "SampleSize": "1200"
        },
        ...
      ]
    }
    ```
- **Error Response:**
  - **Code:** 404
  - **Content:** `{ "error": "No matchups found", "patch": "11.10" }`

### 3. Get All Matchups for a Champion

Retrieves all matchup data for a specific champion in a specific role.

- **URL:** `/matchups/:champion/:role/all`
- **Method:** `GET`
- **URL Parameters:**
  - `champion`: The name of the champion
  - `role`: The role (top, jungle, mid, adc, support)
- **Success Response:**
  - **Code:** 200
  - **Content:** 
    ```json
    {
      "patch": "11.10",
      "matchups": [
        {
          "Champion": "Zed",
          "WinRate": "55.5",
          "SampleSize": "1000"
        },
        {
          "Champion": "Yasuo",
          "WinRate": "52.3",
          "SampleSize": "1200"
        },
        ...
      ]
    }
    ```
- **Error Response:**
  - **Code:** 404
  - **Content:** `{ "error": "No matchups found", "patch": "11.10" }`
