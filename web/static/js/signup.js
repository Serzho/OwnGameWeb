// noinspection ExceptionCaughtLocallyJS

document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                name: '',
                email: '',
                password: '',
                confirmPassword: '',
                isLoading: false,
                errorMessage: '',
                successMessage: ''
            }
        },
        computed: {
            emailValid() {
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                return emailRegex.test(this.email);
            },
            passwordsMatch() {
                return this.password === this.confirmPassword;
            },
            formValid() {
                return this.name.length >= 2 &&
                    this.emailValid &&
                    this.password.length >= 6 &&
                    this.passwordsMatch;
            }
        },
        methods: {
            async handleSubmit() {
                if (!this.formValid) {
                    this.errorMessage = 'Пожалуйста, заполните все поля корректно';
                    return;
                }

                this.isLoading = true;
                this.errorMessage = '';
                this.successMessage = '';

                try {
                    const response = await fetch('/auth/signup', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            name: this.name,
                            email: this.email,
                            password: this.password
                        })
                    });

                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Ошибка регистрации');
                    }

                    this.successMessage = 'Регистрация прошла успешно! Перенаправление...';
                    setTimeout(() => {
                        window.location.href = '/auth/signin';
                    }, 2000);
                } catch (error) {
                    this.errorMessage = error.message;
                } finally {
                    this.isLoading = false;
                }
            }
        }
    }).mount('#app');
});