<template>
    <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">Название склада</label>
                <div class="control">
                    <input v-model="formData.stockName" class="input" type="text">
                </div>
            </div>

            <div class="label"> {{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">добавить склад</button>
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
                stockName: ""
            },
            resp: ""
        };
    },
    methods: {
        async submitModalData() {

            const requestData = {
                StorageName: this.formData.stockName
            }

            try {
                const response = await fetch("http://127.0.0.1:9000/stock/add", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                this.resp = "склад успешно добавлена, ID склада - " + responseData.Data.StorageID

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
