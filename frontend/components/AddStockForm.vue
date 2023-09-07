<template>
    <b-modal v-model="isComponentModalActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
        aria-label="Добавление склада" close-button-aria-label="Закрыть" aria-modal>
        <div class="modal-card">
            <header class="modal-card-head">
                <p class="modal-card-title">Добавление склада</p>
                <button type="button" class="delete" @click="closeModal"></button>
            </header>

            <section class="modal-card-body">
                <b-field label="Название склада">
                    <b-input v-model="formData.stockName" type="text" placeholder="Введите название склада"
                        required></b-input>
                </b-field>
            </section>

            <footer class="modal-card-foot">
                <b-button label="Закрыть" @click="closeModal" />
                <b-button label="Добавить" type="is-primary" @click="submitModalData" />
            </footer>

        </div>
    </b-modal>
</template>
  
<script>
export default {
    data() {
        return {
            isComponentModalActive: this.modalVisible,
            formData: {
                stockName: ''
            },
        }
    },

    props: {
        modalVisible: {
            type: Boolean,
            required: true
        }
    },

    methods: {
        async submitModalData() {
            try {
                const requestData = {
                    storage_name: this.formData.stockName
                };

                const response = await fetch('http://127.0.0.1:9000/stock/add', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });

                const responseData = await response.json();

                this.$emit("fetchStockList")
                this.$buefy.snackbar.open(`склад успешно добавлен, ID склада - ${responseData.Data.stockID}`)
            }
            catch (error) {
                console.error('Ошибка при отправке запроса:', error);
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
            }
            finally {
                this.closeModal()
            }
        },
        
        closeModal() {
            this.$emit("closeModal")
        },
    }
};
</script>