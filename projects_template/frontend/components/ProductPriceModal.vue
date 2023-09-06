<template>
    <section>
        <b-modal v-model="isActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
            aria-label="Добавление цены" close-button-aria-label="Закрыть" aria-modal>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">Добавление цены</p>
                    <button type="button" class="delete" @click="closeModal"></button>
                </header>

                <section class="modal-card-body">
                    <b-field label="Объем продукта">
                        <b-select placeholder="Выберите объем продукта" v-model="formData.variant_id" type="text">
                            <option v-for="opt in variantOptions" :value="opt.Value" :key="opt.Option">{{ opt.weight }}{{
                                opt.unit }}
                            </option>
                        </b-select>
                    </b-field>
                    <b-field label="Цена продукта">
                        <b-input type="number" v-model="formData.price" placeholder="Введите цену" required></b-input>
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
                price: ""
            },
            o: null,
            isActive: this.modalVisible,
            variantOptions: this.options,
        }
    }, props: {
        modalVisible: {
            Type: Boolean,
            required: true
        },

        options: {
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
                price: +this.formData.price
            }

            try {
                const response = await fetch("http://127.0.0.1:9000/product/price", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                })

                const responseData = await response.json()

                this.$buefy.snackbar.open(`Цена успешно добавлена, priceID - ${responseData.Data.price_id}`)
            } catch (error) {
                this.$buefy.snackbar.open(`${error}`)
                console.error(error)
            } finally {
                this.closeModal()
            }
        }
    }
}
</script>