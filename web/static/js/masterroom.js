document.addEventListener('DOMContentLoaded', () => {
    const {createApp} = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                roomName: '',
                themes: [],
                players: [],
                maxPlayers: 0,
                currentPlayers: 0,
                activeQuestion: null,
                usedQuestions: [],
                answerPhase: false,
                respondingPlayer: null,
                timer: 10,
                answerInterval: null,
                checkAnswerInterval: null,
                showAnswerButtons: false,
                errorMessage: ''
            }
        },
        mounted() {
            this.fetchGameData();
            this.fetchQuestions();
        },
        methods: {
            async fetchGameData() {
                try {
                    const response = await fetch('/play/gameinfo');
                    const data = await response.json();
                    const jdata = JSON.parse(data.data);
                    this.roomName = jdata.title;
                    this.players = jdata.players.map(p => ({...p, score: 0}));
                    this.maxPlayers = jdata.maxPlayers;
                    this.currentPlayers = jdata.players.length;
                } catch (error) {
                    this.errorMessage = 'Ошибка загрузки данных игры';
                }
            },

            async fetchQuestions() {
                try {
                    const response = await fetch('/play/questions');
                    const data = await response.json();
                    const jdata = JSON.parse(data.data);
                    this.themes = jdata.first_round;
                } catch (error) {
                    this.errorMessage = 'Ошибка загрузки вопросов';
                }
            },

            selectQuestion(question) {
                this.activeQuestion = question;
                this.answerPhase = false;
                this.respondingPlayer = null;
            },

            async startAnswerPhase() {
                // Запуск фазы ответа
                this.answerPhase = true;
                const delay = Math.random() * 700 + 100;

                setTimeout(async () => {
                    try {
                        await fetch('/play/startanswer', {method: 'POST'});
                        this.startAnswerPolling();
                    } catch (error) {
                        this.errorMessage = error.message;
                    }
                }, delay);
            },

            startAnswerPolling() {
                this.checkAnswerInterval = setInterval(async () => {
                    const response = await fetch('/play/checkanswer');
                    const player = await response.json();

                    if (player) {
                        this.respondingPlayer = player;
                        this.startResponseTimer();
                        clearInterval(this.checkAnswerInterval);
                    }
                }, 500);

                // Таймер ожидания ответа
                setTimeout(() => {
                    if (!this.respondingPlayer) {
                        this.closeQuestion();
                        clearInterval(this.checkAnswerInterval);
                    }
                }, 10000);
            },

            startResponseTimer() {
                this.timer = 5;
                this.answerInterval = setInterval(() => {
                    this.timer--;
                    if (this.timer <= 0) {
                        this.showAnswerButtons = true;
                        clearInterval(this.answerInterval);
                    }
                }, 1000);
            },

            async handleAnswer(isCorrect) {
                const answer = {
                    playerId: this.respondingPlayer.id,
                    questionId: this.activeQuestion.question_id,
                    correct: isCorrect
                };

                try {
                    await fetch('/play/answer', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify(answer)
                    });

                    if (isCorrect) {
                        this.players.find(p => p.id === answer.playerId).score += this.activeQuestion.level;
                        this.closeQuestion();
                    } else {
                        this.respondingPlayer = null;
                        this.startAnswerPhase();
                    }
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            closeQuestion() {
                this.usedQuestions.push(this.activeQuestion.question_id);
                this.activeQuestion = null;
                this.answerPhase = false;
                this.respondingPlayer = null;
                this.showAnswerButtons = false;
                clearInterval(this.checkAnswerInterval);
                clearInterval(this.answerInterval);
            }
        }
    }).mount('#app');
});
