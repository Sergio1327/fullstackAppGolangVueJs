<template>
    <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">ID варианта продукта</label>
                <div class="control">
                    <input v-model="formData.variant_id" class="input" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">ID склада</label>
                <div class="control">
                    <input v-model="formData.storage_id" class="input" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">Кол-во продуктов</label>
                <div class="control">
                    <input v-model="formData.quantity" class="input" type="number">
                </div>
            </div>
            <div class="label"> {{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">добавить продажу</button>
            <button class="button" @click="closeModal">Закрыть</button>
        </div>
        <button class="modal-close is-large" aria-label="Закрыть" @click="closeModal"></button>

    </div>
</template>

<script>
export default {
    data() {
        return {
            formData: {
                variant_id: 1,
                storage_id: 3,
                quantity: 2
            },
            resp: ""
        };
    },
    methods: {
        async submitModalData() {



            const requestData = {
                variant_id: parseInt(this.formData.variant_id),
                storage_id: parseInt(this.formData.storage_id),
                quantity: parseInt(this.formData.quantity)
            }

            try {
                const response = await fetch("http://127.0.0.1:9000/buy", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                this.resp = "Продажа успешно добавлена, ID продажи - " + responseData.Data.sale_id

            } catch (error) {
                console.error("Ошибка при отправке запроса:", error);
                console.log(error)

                this.resp = error
            }

        },
        closeModal() {
            this.$emit('closeModal');
        },
    },
}
</script>
