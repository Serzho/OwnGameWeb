document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                code: '',
                isLoading: false,
                errorMessage: '',
                successMessage: ''
            }
        },
        computed: {
            formValid() {
                return /^[A-Za-z0-9]{6}$/.test(this.code);
            }
        },
        methods: {
            async handleSubmit() {
                if (!this.formValid) {
                    this.errorMessage = 'Код должен содержать 6 символов (буквы/цифры)';
                    return;
                }

                this.isLoading = true;
                this.errorMessage = '';
                this.successMessage = '';

                try {
                    const response = await fetch('/joingame', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            code: this.code
                        })
                    });

                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Ошибка присоединения');
                    }

                    this.successMessage = 'Успешное подключение! Перенаправление...';
                    setTimeout(() => {
                        window.location.href = '/play/waitingroom';
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