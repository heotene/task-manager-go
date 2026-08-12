document.addEventListener("DOMContentLoaded", function () {

    const toggle = document.getElementById("darkModeToggle");

    // Load saved theme
    if (localStorage.getItem("theme") === "dark") {
        document.documentElement.classList.add("dark");
        
        if (toggle) {
            toggle.checked = true;
        }
    }

    // Theme switch
    if (toggle) {

        toggle.addEventListener("change", function () {

            if (this.checked) {

                document.documentElement.classList.add("dark");

                localStorage.setItem("theme", "dark");

            } else {

                document.documentElement.classList.remove("dark");

                localStorage.setItem("theme", "light");

            }

        });

    }

});

// Confirm account deletion
document.addEventListener("DOMContentLoaded", function () {

    const deleteForm = document.querySelector(
        'form[action="/settings/delete-account"]'
    );

    if (deleteForm) {

        deleteForm.addEventListener("submit", function (event) {

            const confirmed = confirm(
                "⚠️ Are you sure you want to permanently delete your account?\n\n" +
                "This will delete your account and all your tasks.\n\n" +
                "This action cannot be undone."
            );

            if (!confirmed) {
                event.preventDefault();
            }

        });

    }

});
