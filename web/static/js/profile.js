document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                playedGames: 0,
                wonGames: 0,
                currentName: '',
                newName: '',
                oldPassword: '',
                newPassword: '',
                errorMessage: '',
                successMessage: ''
            }
        },
        computed: {
            formValid() {
                return this.oldPassword.length > 0 &&
                    (this.newName.length > 0 || this.newPassword.length >= 6);
            }
        },
        methods: {
            async loadProfile() {
                try {
                    const response = await fetch(`/profile/info`);
                    const data = await response.json();

                    if (!response.ok) throw new Error(data.message || 'Ошибка загрузки');
                    const jdata = JSON.parse(data.data);
                    this.playedGames = jdata.playedGames;
                    this.wonGames = jdata.wonGames;
                    this.currentName = jdata.name;
                } catch (error) {
                    this.errorMessage = 'Ошибка загрузки данных профиля';
                }
            },
            async submitForm() {
                this.errorMessage = '';
                this.successMessage = '';

                const formData = {
                    oldPassword: this.oldPassword,
                    newName: this.newName,
                    newPassword: this.newPassword
                };

                try {
                    const response = await fetch('/profile/update', {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify(formData)
                    });

                    const result = await response.json();

                    if (response.ok) {
                        this.successMessage = 'Данные успешно обновлены!';
                        if (result.new_name) this.currentName = result.new_name;
                        this.newName = '';
                        this.newPassword = '';
                        this.oldPassword = '';
                    } else {
                        this.errorMessage = result.error || 'Ошибка обновления данных';
                    }
                } catch (error) {
                    this.errorMessage = 'Ошибка соединения';
                }
            }
        },
        mounted() {
            this.loadProfile();
        }
    }).mount('#app');
});