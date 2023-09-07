<template>
    <b-modal v-model="isComponentModalActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
        aria-label="Добавление склада" close-button-aria-label="Закрыть" aria-modal>
        <div class="modal-card">
            <header class="modal-card-head">
                <p class="modal-card-title">Добавление продукта</p>
                <button type="button" class="delete" @click="closeModal"></button>
            </header>

            <section class="modal-card-body">
                <b-field label="Название продукта">
                    <b-input v-model="formData.productName" type="text" placeholder="Введите название продукта"
                        required></b-input>
                </b-field>
                <b-field label="Описание продукта">
                    <b-input v-model="formData.descr" type="text" placeholder="Введите описание продукта"
                        required></b-input>
                </b-field>
                <b-field label="Теги продукта">
                    <b-input v-model="formData.tags" type="text" placeholder="Введите теги продукта" required></b-input>
                </b-field>
                <b-field class="" label="Варианты продукта">
                    <div v-for="(v, index) in formData.productVariants" :key="index">
                        <b-input v-model="v.weight" type="number" placeholder="Объем" required></b-input>
                        <b-input v-model="v.unit" class="mt-1 mb-5" type="text" placeholder="Eдиница измерения"
                            required></b-input>
                    </div>
                </b-field>
                <b-field class="is-flex is-justify-content-flex-end">
                    <b-button label="добавить вариант" @click="addVariant"></b-button>
                </b-field>
                <b-field class="is-flex is-justify-content-flex-end">
                    <b-button label="удалить вариант" @click="deleteVariant"></b-button>
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
                productName: "",
                descr: "",
                addedAt: null,
                tags: "",
                productVariants: [{
                    weight: "",
                    unit: ""
                }]
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
        addVariant() {
            this.formData.productVariants.push({})
        },
        deleteVariant() {
            if (this.formData.productVariants.length > 0) {
                this.formData.productVariants.pop()
            }
        },

        async submitModalData() {
            const formatedVariants = this.formData.productVariants.map(v => ({
                weight: +v.weight,
                unit: v.unit
            }))

            const requestData = {
                name: this.formData.productName,
                description: this.formData.descr,
                tags: this.formData.tags,
                variants: formatedVariants
            }

            try {
                const response = await fetch("http://127.0.0.1:9000/product/add", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                this.$buefy.snackbar.open(`Продук успешно добавлен в базу, ID продукта - ${responseData.Data.product_id}`)

            }
            catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error("Ошибка при отправке запроса:", error);
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