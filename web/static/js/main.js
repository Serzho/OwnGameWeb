document.addEventListener('DOMContentLoaded', () => {
    const { createApp } = Vue;

    createApp({
        compilerOptions: {
            delimiters: ['[[', ']]']
        },
        methods: {
            async signOut() {
                try {
                    const response = await fetch(`/auth/signout`, {
                        method: 'POST'
                    });

                    window.location.href = '/'
                } catch (error) {
                    this.errorMessage = error.message;
                }
            },
        }
    }).mount('#app');
});