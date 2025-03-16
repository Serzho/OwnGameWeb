document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                roomName: 'Загрузка...',
                players: [],
                maxPlayers: 4,
                currentPlayers: 0,
                isHost: false,
                gameId: null,
                isLoading: false,
                errorMessage: '',
                updateInterval: null
            }
        },
        computed: {
            canStart() {
                return this.currentPlayers >= 2 && !this.isLoading && this.isHost;
            }
        },
        mounted() {
            this.gameId = window.location.pathname.split('/').pop();
            this.loadData();
            this.updateInterval = setInterval(this.loadData, 10000);
        },
        beforeUnmount() {
            clearInterval(this.updateInterval);
        },
        methods: {
            async loadData() {
                try {
                    const response = await fetch(`/play/gameinfo`);
                    const data = await response.json();

                    if (!response.ok) throw new Error(data.message || 'Ошибка загрузки');

                    this.roomName = data.data.title;
                    this.players = data.data.players;
                    this.maxPlayers = data.data.maxPlayers;
                    this.currentPlayers = data.data.players.length;
                    this.isHost = data.data.isHost;
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            async removePlayer(playerId) {
                try {
                    const response = await fetch(`/play/removeplayer`, {
                        method: 'DELETE',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ playerId })
                    });

                    if (!response.ok) throw new Error('Ошибка удаления игрока');
                    await this.loadData();
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            async cancelGame() {
                if (!confirm('Вы уверены, что хотите отменить игру?')) return;

                try {
                    const response = await fetch(`/play/game`, {
                        method: 'DELETE'
                    });

                    if (!response.ok) throw new Error('Ошибка отмены игры');
                    window.location.href = '/main';
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            async startGame() {
                this.isLoading = true;
                this.errorMessage = '';

                try {
                    const response = await fetch(`/play/start`, {
                        method: 'POST'
                    });

                    if (!response.ok) throw new Error('Ошибка запуска игры');
                    window.location.href = `/game/${this.gameId}/play`;
                } catch (error) {
                    this.errorMessage = error.message;
                } finally {
                    this.isLoading = false;
                }
            }
        }
    }).mount('#app');
});