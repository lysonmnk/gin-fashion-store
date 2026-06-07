// BUG FIX: Sebelumnya menggunakan fetch('/logout') yang tidak mengikuti
// server-side redirect dengan benar. Cukup arahkan browser langsung ke /logout.
function handleLogout() {
    window.location.href = "/logout";
}