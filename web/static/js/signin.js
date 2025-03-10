// noinspection ExceptionCaughtLocallyJS

document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                email: '',
                password: '',
                isLoading: false,
                errorMessage: ''
            }
        },
        computed: {
            emailValid() {
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                return emailRegex.test(this.email);
            },
            formValid() {
                return this.emailValid && this.password.length >= 6;
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

                try {
                    const response = await fetch('/auth/signin', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            email: this.email,
                            password: this.password
                        })
                    });

                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Ошибка авторизации');
                    }

                    window.location.href = '/main';
                } catch (error) {
                    this.errorMessage = error.message;
                } finally {
                    this.isLoading = false;
                }
            }
        }
    }).mount('#app');
});