// noinspection ExceptionCaughtLocallyJS

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
                errorMessage: '',
                showEditModal: false,
                editingPack: { id: null, title: '', content: '' },
                showServerPacks: false,
                serverPacks: [],
                searchQuery: ''
            }
        },
        computed: {
            filteredServerPacks() {
                return this.serverPacks.filter(pack =>
                    pack.title.toLowerCase().includes(this.searchQuery.toLowerCase())
                );
            }
        },
        mounted() {
            this.loadPacks();
        },
        methods: {

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

            async signOut(packId) {
                try {
                    const response = await fetch(`/auth/signout`, {
                        method: 'POST'
                    });

                    window.location.href = '/'
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },

            navigateProfile() {
                window.location.href = '/profile';
            },
            async editPack(packId) {
                try {
                    const response = await fetch(`/getpack/${packId}`);
                    const data = await response.json();
                    this.editingPack = {
                        id: packId,
                        title: data.title,
                        content: data.content
                    };
                    this.showEditModal = true;
                } catch (error) {
                    this.errorMessage = 'Ошибка загрузки пакета';
                }
            },

            async savePackName() {
                try {
                    await fetch(`/updatepacktitle/${this.editingPack.id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ title: this.editingPack.title })
                    });
                    await this.loadPacks();
                } catch (error) {
                    this.errorMessage = 'Ошибка сохранения имени';
                }
            },

            async savePackFile() {
                try {
                    await fetch(`/updatepackfile/${this.editingPack.id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'text/csv' },
                        body: JSON.stringify({ content: this.editingPack.content})
                    });
                    this.showEditModal = false;
                } catch (error) {
                    this.errorMessage = 'Ошибка сохранения файла';
                }
            },

            closeModal() {
                this.showEditModal = false;
                this.editingPack = { id: null, title: '', content: '' };
            },

            async loadServerPacks() {
                try {
                    const response = await fetch('/getserverpacks');
                    const data = await response.json();
                    this.serverPacks = data.packs;
                    this.showServerPacks = true;
                } catch (error) {
                    this.errorMessage = 'Ошибка загрузки пакетов';
                }
            },

            async addServerPack(packId) {
                try {
                    await fetch(`/addserverpack/${packId}`, { method: 'POST' });
                    await this.loadPacks();
                    this.showServerPacks = false;
                } catch (error) {
                    this.errorMessage = 'Ошибка добавления пакета';
                }
            }
        }
    }).mount('#app');
});