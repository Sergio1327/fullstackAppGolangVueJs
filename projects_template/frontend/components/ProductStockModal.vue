<template>
    <section>
        <b-modal v-model="isActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
            aria-label="Добавление продукта на склад" close-button-aria-label="Закрыть" aria-modal>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">Добавление цены</p>
                    <button type="button" class="delete" @click="closeModal"></button>
                </header>

                <section class="modal-card-body">
                    <b-field label="Объем продукта">
                        <b-select v-model="formData.variant_id" type="text" placeholder="Выберите объем продукта" required>
                            <option v-for="opt in variantOptions" :value="opt.Value" :key="opt.Option">{{ opt.weight }}{{
                                opt.unit }}
                            </option>
                        </b-select>
                    </b-field>
                    <b-field label="Название склада">
                        <b-select v-model="formData.storage_id" type="text" placeholder="Выберите склад" required>
                            <option v-for="s in stockOptions" :value="s.Value" :key="s.Option">{{ s.StorageName }}
                            </option>
                        </b-select>
                    </b-field>
                    <b-field label="Колличество">
                        <b-input type="number" v-model="formData.quantity" placeholder="Введите колличество продукта"
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
            formData: {
                variant_id: null,
                storage_id: null,
                quantity: ""
            },

            isActive: this.modalVisible,
            variantOptions: this.options,
            stockOptions: this.storageOptions,

        }
    }, props: {
        modalVisible: {
            Type: Boolean,
            required: true
        },

        options: {
            Type: Array,
            required: true
        },

        storageOptions: {
            Type: Array,
            required: true
        }
    },
    methods: {
        closeModal() {
            this.$emit("closeModal")
        },
        async submitModalData() {

            const requestData = {
                variant_id: +this.formData.variant_id,
                storage_id: +this.formData.storage_id,
                quantity: +this.formData.quantity
            }

            try {
                const response = await fetch("http://127.0.0.1:9000/product/add/stock", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                })

                const responseData = await response.json()
                this.$buefy.snackbar.open(`Продукт успешно добавлен на склад,ID операции  - ${responseData.Data.product_stock_ID}`)
            } catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type:"is-danger"
                })
                console.error(error)
            } finally {
                this.closeModal()
            }
        }
    }
}
</script>