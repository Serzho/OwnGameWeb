// noinspection ExceptionCaughtLocallyJS

document.addEventListener('DOMContentLoaded', () => {
    const {createApp} = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                email: '',
                isLoading: false,
                errorMessage: '',
                successMessage: ''
            }
        },
        computed: {
            emailValid() {
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                return emailRegex.test(this.email);
            }
        },
        methods: {
            async handleSubmit() {
                if (!this.emailValid) {
                    this.errorMessage = 'Пожалуйста, введите корректный email';
                    return;
                }

                this.isLoading = true;
                this.errorMessage = '';
                this.successMessage = '';

                try {
                    const response = await fetch('/auth/recoverPassword', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ email: this.email })
                    });

                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Ошибка сервера');
                    }

                    this.successMessage = 'Инструкции по восстановлению отправлены на ваш email';
                    this.email = '';
                } catch (error) {
                    this.errorMessage = error.message;
                } finally {
                    this.isLoading = false;
                }
            }
        }
    }).mount('#app');
});