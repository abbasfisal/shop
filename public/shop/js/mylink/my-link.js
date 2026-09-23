// پیدا کردن تمام لینک‌ها
const links = document.querySelectorAll('.my-link');

// اضافه کردن رویداد کلیک به هر لینک
links.forEach(link => {
    link.addEventListener('click', function (event) {
        // اعمال کلاس active به لینک کلیک‌شده
        links.forEach(link => link.classList.remove('active')); // حذف کلاس active از بقیه لینک‌ها
        this.classList.add('active'); // اضافه کردن کلاس active به لینک کلیک‌شده

        // ریدایرکت به آدرس مقصد (بدون جلوگیری از ریدایرکت)
        window.location.href = this.href;
    });
});
