// Operasi Fetch API Keranjang Belanja
async function apiAddToCart(productId, quantity) {
    return await fetch('/api/cart', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ product_id: productId, quantity: quantity })
    });
}

async function apiUpdateCart(productId, quantity) {
    return await fetch('/api/cart', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ product_id: productId, quantity: quantity })
    });
}

async function apiRemoveFromCart(productId) {
    return await fetch('/api/cart', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ product_id: productId })
    });
} 