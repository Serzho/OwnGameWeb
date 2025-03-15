document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        data() {
            return {
                packs: [],
                isLoading: false,
                errorMessage: ''
            }
        },
        mounted() {
            this.loadPacks();
        },
        methods: {
            addFromServer() {
                // Заглушка для будущей реализации
                console.log('Добавление с сервера');
            },

            downloadPack(packId) {
                window.location.href = `/downloadpack/${packId}`;
            },

            async loadPacks() {
                try {
                    const response = await fetch('/getallpacks');
                    if (!response.ok) throw new Error('Ошибка загрузки');
                    const data = await response.json();
                    this.packs = data.packs;
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            async handleFileUpload(event) {
                const file = event.target.files[0];
                if (!file) return;

                try {
                    const formData = new FormData();
                    formData.append('packFile', file);

                    const response = await fetch('/addpack', {
                        method: 'POST',
                        body: formData
                    });

                    if (!response.ok) throw new Error('Ошибка загрузки');
                    await this.loadPacks();
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            async deletePack(packId) {
                if (!confirm('Удалить пакет?')) return;

                try {
                    const response = await fetch(`/deletepack/${packId}`, {
                        method: 'DELETE'
                    });

                    if (!response.ok) throw new Error('Ошибка удаления');
                    this.packs = this.packs.filter(p => p.id !== packId);
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            editPack(packId) {
                window.location.href = `/editpack/${packId}`;
            },

            navigateProfile() {
                window.location.href = '/profile';
            }
        }
    }).mount('#app');
});