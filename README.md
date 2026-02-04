# 🚁 Drone Delivery Management Backend

A JWT-secured REST API for managing drone deliveries.

---

## 🔐 Authentication
- Generate token via: `POST /login`
- Add token to **Postman Collection Authorization**
- Set requests to **Inherit auth from parent**
- No need to manually add `Bearer`

---

## 👥 Roles
- `admin`
- `customer`
- `drone`

---

## 🚁 Drone Status
- `Available`
- `Busy`
- `Broken`
- `Fixed`

---

## 📦 Order Status
- `Pending`
- `Assigned`
- `Failed`
- `Delivered`

---

## 🧪 Testing
1. Call `/login`
2. Copy JWT token
3. Add it to Postman Collection Authorization
4. Test all endpoints

---

## ⚙️ Features

### Drones
- Reserve jobs  
- Pickup orders  
- Deliver or fail orders  
- Update location  
- Mark broken/fixed  
- Get current assigned order  

### Customers
- Create delivery orders  
- Cancel pending orders  
- Track order status  

### Admins
- View all orders  
- Update origin/destination  
- Manage drones  
- Mark drones broken/fixed  

---

## 📌 Rules
- If a drone becomes **Broken**, its order is automatically reassigned to another available drone.
