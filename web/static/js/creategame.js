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
                csvFile: null,
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
                    this.csvFile !== null
            }
        },
        methods: {
            handleFileUpload(event) {
                const file = event.target.files[0];
                if(file && file.name.endsWith('.csv')) {
                    this.csvFile = file;
                } else {
                    this.errorMessage = 'Пожалуйста, выберите CSV файл';
                    this.csvFile = null;
                }
            },

            async handleSubmit() {
                if(!this.formValid) {
                    this.errorMessage = 'Заполните все поля корректно';
                    return;
                }

                this.isLoading = true;
                this.errorMessage = '';
                this.successMessage = '';

                try {
                    const formData = new FormData();
                    formData.append('roomName', this.roomName);
                    formData.append('maxPlayers', this.maxPlayers);
                    formData.append('questionsFile', this.csvFile);

                    const response = await fetch('/game/create', {
                        method: 'POST',
                        body: formData
                    });

                    const data = await response.json();

                    if(!response.ok) {
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