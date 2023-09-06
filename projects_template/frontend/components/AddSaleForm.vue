<template>
    <section>
        <b-modal v-model="isComponentModalActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
            aria-label="Добавление склада" close-button-aria-label="Закрыть" aria-modal>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">Добавление продажи</p>
                    <button type="button" class="delete" @click="closeModal"></button>
                </header>

                <section class="modal-card-body">
                    <b-field label="Вариант продукта">
                        <b-select v-model="formData.variant_id" type="text" placeholder="Выберите вариант продукта"
                            required>
                            <option v-for="v in vOptions " :value="v.Option" :key="v.Value">{{ v.ProductName }} {{ v.Weight
                            }}{{ v.Unit }}</option>
                        </b-select>
                    </b-field>
                    <b-field label="Выберите склад">
                        <b-select v-model="formData.storage_id" type="text" placeholder="Выберите склад" required>
                            <option v-for="s in stockOptions" :value="s.Option" :key="s.Value">{{ s.StorageName }}</option>
                        </b-select>
                    </b-field>
                    <b-field label="Введите колличество">
                        <b-input v-model="formData.quantity" type="text" placeholder="Введите колличество"
                            required></b-input>
                    </b-field>
                </section>

                <footer class="modal-card-foot">
                    <b-button label="Закрыть" @click="closeModal" />
                    <b-button label="Добавить" type="is-primary" @click="submitModalData" />
                </footer>

            </div>
        </b-modal>
    </section>
</template>
  
<script>
export default {
    data() {
        return {
            isComponentModalActive: this.modalVisible,
            vOptions: this.variantOptions,

            formData: {
                variant_id: null,
                storage_id: null,
                quantity: ""
            },


        }
    },
    props: {
        modalVisible: {
            type: Boolean,
            required: true
        },

        variantOptions: {
            type: Array,
            required: true
        },

        stockOptions: {
            type: Array,
            required: true
        }
    },
    methods: {
        async submitModalData() {
            try {
                const requestData = {
                    variant_id: +this.formData.variant_id,
                    storage_id: +this.formData.storage_id,
                    quantity: +this.formData.quantity
                };

                const response = await fetch('http://127.0.0.1:9000/buy', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });

                const responseData = await response.json();

                this.$buefy.snackbar.open(`продажа успешно добавлена, ID продажи - ${responseData.Data.sale_id}`)
            } catch (error) {
                console.error('Ошибка при отправке запроса:', error);
                this.$buefy.snackbar.open(error)
            } finally {
                this.closeModal()
            }
        },
        closeModal() {
            this.vOptions = []
            this.$emit("closeModal")

        }
    }
};
</script>