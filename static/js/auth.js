// Fungsi utilitas autentikasi global
async function handleLogout() {
    try {
        const response = await fetch('/logout');
        if (response.ok) {
            window.location.href = "/login";
        }
    } catch (err) {
        console.error("Gagal memproses keluar sesi:", err);
    }
}