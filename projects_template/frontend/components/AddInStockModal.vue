<template>
    <div class="modal" :class="{ 'is-active': showStockModal }">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">ID варианта продукта</label>
                <div class="control">
                    <input class="input" v-model="req.variant_id" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">ID склада</label>
                <div class="control">
                    <input class="input" v-model="req.storage_id" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">Дата добавления продукта на склад</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="req.added_at">
                </div>
            </div>
            <div class="field">
                <label class="label">Кол-во</label>
                <div class="control">
                    <input class="input" type="number" v-model="req.quantity">
                </div>
            </div>

            <div class="label">{{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">Добавить цену продукта</button>
            <button class="button" @click="closeModal">Закрыть</button>
        </div>
        <button class="modal-close is-large" aria-label="Закрыть" @click="closeModal"></button>
    </div>
</template>

<script  >
export default {
    props: {
        showStockModal: {
            type: Boolean,
            required: true
        }
    },
    data() {
        return {
            req: {
                variant_id: "",
                added_at: null,
                storage_id: "",
                quantity: ""

            },
            resp: ""
        }
    },
    methods: {
        closeModal() {
            this.$emit("closeModal")
        },
        formatDate(dateTime) {
            const date = new Date(dateTime);
            const year = date.getUTCFullYear();
            const month = String(date.getUTCMonth() + 1).padStart(2, "0");
            const day = String(date.getUTCDate()).padStart(2, "0");
            const hours = String(date.getHours()).padStart(2, "0");
            const minutes = String(date.getUTCMinutes()).padStart(2, "0");
            const seconds = String(date.getUTCSeconds()).padStart(2, "0");

            return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}+05:00`;
        },
        async submitModalData() {
            try {

                const requestData = {
                    variant_id: parseInt(this.req.variant_id),
                    storage_id: parseInt(this.req.storage_id),
                    added_at: this.formatDate(this.req.added_at),
                    quantity: parseInt(this.req.quantity)
                }

                const response = await fetch("http://127.0.0.1:9000/product/add/stock", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                this.resp = `Продукт успешно добавлен на склад, ID операции = ${responseData.Data.product_stock_ID}`
            } catch (error) {
                console.error(error)
                this.resp = error
            }
        },
    }
}
</script>