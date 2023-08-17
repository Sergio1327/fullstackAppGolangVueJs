<template>
    <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">Variant ID</label>
                <div class="control">
                    <input class="input" type="number" v-model="formData.variant_id">
                </div>
            </div>
            <div class="field">
                <label class="label">Storage ID</label>
                <div class="control">
                    <input class="input" type="number" v-model="formData.storage_id">
                </div>
            </div>
            <div class="field">
                <label class="label">Quantity</label>
                <div class="control">
                    <input class="input" type="number" v-model="formData.quantity">
                </div>
            </div>
            <button @click="submitModalData" class="button is-primary">Создать запрос</button>
            <button @click="closeModal" class="button">Закрыть</button>
        </div>
        <button @click="closeModal" class="modal-close is-large" aria-label="Закрыть"></button>
        <div class="resp">{{ resp }}</div>
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
                console.log(responseData)

            } catch (error) {
                console.error("Ошибка при отправке запроса:", error);
                console.log(error)
                this.resp = error
            }



            this.closeModal();
        },
        closeModal() {
            this.$emit('closeModal');
        },
    },
}
</script>
