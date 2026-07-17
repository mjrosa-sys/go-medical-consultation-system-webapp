function switchTab(tabName) {
    // Hide all forms
    document.querySelectorAll('.auth-form').forEach(form => {
        form.classList.remove('active');
    });

    // Deactivate all tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });

    // Show selected form and activate button
    if (tabName === 'login') {
        document.getElementById('login-form').classList.add('active');
        event.target.classList.add('active');
    } else if (tabName === 'register') {
        document.getElementById('register-form').classList.add('active');
        event.target.classList.add('active');
    }
}