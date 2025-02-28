const axios = require('axios');

let cart = [];

function fetchMenu() {
    console.log("Fetching menu...");
    
    axios.get("http://localhost:3000/menu")
        .then(response => {
            console.log("Menu data received:", response.data);

            let menuList = document.getElementById("menu-list");
            menuList.innerHTML = "";

            response.data.forEach(menu => {
                let div = document.createElement("div");
                div.className = "menu-card";
                div.innerHTML = `
                    <img src="./assets/${menu.id}.jpg" alt="${menu.name}">
                    <h3>${menu.name}</h3>
                    <p>Rp ${menu.price}</p>
                    <button onclick="addToCart(${menu.id}, '${menu.name}', ${menu.price})">Beli</button>
                `;
                menuList.appendChild(div);
            });
        })
        .catch(error => {
            console.error("Gagal mengambil menu:", error);
            alert("Gagal memuat menu. Pastikan backend berjalan!");
        });
}

function addToCart(menuId, name, price) {
    let item = cart.find(i => i.menu_id === menuId);
    if (item) {
        item.qty++;
    } else {
        cart.push({ menu_id: menuId, name: name, price: price, qty: 1 });
    }
    updateCartUI();
}

function updateCartUI() {
    let cartList = document.getElementById("cart");
    let totalPrice = 0;
    cartList.innerHTML = "";

    cart.forEach((item, index) => {
        let li = document.createElement("li");
        li.className = "cart-item";
        li.innerHTML = `
            <span>${item.name} x${item.qty} - Rp ${item.price * item.qty}</span>
            <button onclick="removeFromCart(${index})">Hapus</button>
        `;
        cartList.appendChild(li);
        totalPrice += item.price * item.qty;
    });

    document.getElementById("total-price").innerText = totalPrice;
}

// Panggil fetchMenu() saat aplikasi dibuka
fetchMenu();
