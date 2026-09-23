//---------------------- for header price
// تابع برای قالب‌بندی قیمت‌ها
function formatPriceElements(selector) {
    const elements = document.querySelectorAll(selector);
    elements.forEach(el => {
        const numText = el.textContent.replace(/[^\d]/g, ''); // حذف تمام کاراکترهای غیر عدد
        if (numText) {
            const formatted = Number(numText).toLocaleString('en-US'); // قالب‌بندی عدد
            el.textContent = formatted + ' تومان'; // اضافه کردن تومان
        }
    });
}

// صبر کردن تا بارگذاری کامل صفحه
window.addEventListener('DOMContentLoaded', () => {
    // قالب‌بندی قیمت‌ها در بخش‌های مختلف
    formatPriceElements('.price');
    formatPriceElements('.total');

    // قیمت‌ها در بخش‌های خاص
    formatPriceElements('.item-price');
    formatPriceElements('.summery-total-original-price');
    formatPriceElements('.summary-discount-total-profile-price');
    formatPriceElements('.checkout-summary-price-value-amount');
});

//------------------------------------ end for header price