// Pengendali filter katalog
function filterByCategory(slug) {
    if (slug === 'all') {
        window.location.href = '/';
    } else {
        window.location.href = `/?category=${slug}`;
    }
}