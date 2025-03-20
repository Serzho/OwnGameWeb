document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                roomName: 'Название комнаты',
                players: [],
                isButtonDisabled: false,
                timer: 0,
                isShaking: false,
                errorMessage: ''
            }
        },
        mounted() {
            this.fetchGameData();
            setInterval(this.fetchGameData, 3000);
        },
        methods: {
            async fetchGameData() {
                try {
                    const response = await fetch('/play/gameinfo');
                    const data = await response.json();
                    const jdata = JSON.parse(data.data);
                    this.roomName = jdata.title;
                    this.players = jdata.players;
                } catch (error) {
                    this.errorMessage = 'Ошибка обновления данных';
                }
            },

            async handleAnswer() {
                try {
                    const response = await fetch('/play/buzz', { method: 'POST' });
                    const result = await response.json();

                    switch(result.data) {
                        case 'falsestart':
                            this.handleFalseStart();
                            break;
                        case 'in time':
                            this.handleInTime();
                            break;
                    }
                } catch (error) {
                    this.errorMessage = 'Ошибка отправки ответа';
                }
            },

            handleFalseStart() {
                this.isButtonDisabled = true;
                this.isShaking = true;

                setTimeout(() => {
                    this.isButtonDisabled = false;
                    this.isShaking = false;
                }, 1000);
            },

            handleInTime() {
                this.isButtonDisabled = true;
                this.timer = 10;

                const interval = setInterval(() => {
                    this.timer--;
                    if(this.timer <= 0) {
                        clearInterval(interval);
                        this.isButtonDisabled = false;
                    }
                }, 1000);
            }
        }
    }).mount('#app');
});