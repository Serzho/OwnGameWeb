document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                roomName: '',
                maxPlayers: 4,
                selectedPackId: null,
                packs: [],
                isLoading: false,
                errorMessage: '',
                successMessage: ''
            }
        },
        computed: {
            formValid() {
                return this.roomName.length > 0 &&
                    this.maxPlayers >= 2 &&
                    this.maxPlayers <= 6 &&
                    this.selectedPackId !== null
            }
        },
        mounted() {
            this.loadPacks();
        },
        methods: {
            async loadPacks() {
                try {
                    const response = await fetch('/getallpacks');
                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Не удалось загрузить паки');
                    }

                    this.packs = data.packs;
                    if (this.packs.length === 0) {
                        this.errorMessage = 'Нет доступных паков вопросов';
                    }
                } catch (error) {
                    this.errorMessage = `Ошибка: ${error.message}`;
                    console.error('Ошибка загрузки паков:', error);
                }
            },

            async handleSubmit() {
                if (!this.formValid) {
                    this.errorMessage = 'Заполните все поля корректно';
                    return;
                }

                this.isLoading = true;
                this.errorMessage = '';
                this.successMessage = '';

                try {
                        const response = await fetch('/creategame', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            title: this.roomName,
                            maxPlayers: this.maxPlayers,
                            packId: this.selectedPackId
                        })
                    });

                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.message || 'Ошибка при создании игры');
                    }

                    this.successMessage = 'Игра успешно создана! Перенаправление...';
                    setTimeout(() => {
                        window.location.href = `/game/${data.gameId}`;
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