<template>
    <div class="modal" :class="{ 'is-active': showPriceModal }">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">ID варианта продукта</label>
                <div class="control">
                    <input class="input" v-model="req.variant_id" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">Дата начала цены</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="req.start_date">
                </div>
            </div>
            <div class="field">
                <label class="label">Дата окончания цены</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="req.end_date">
                </div>
            </div>
            <div class="field">
                <label class="label">Цена продукта</label>
                <div class="control">
                    <input class="input" type="number" v-model="req.price">
                </div>
            </div>

            <div class="label">{{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">Добавить цену продукта</button>
            <button class="button" @click="closeModal">Закрыть</button>
        </div>
        <button class="modal-close is-large" aria-label="Закрыть" @click="closeModal"></button>
    </div>
</template>

<script>
export default {
    props: {
        showPriceModal: {
            type: Boolean,
            required: true
        }
    },
    data() {
        return {
            req: {
                variant_id: "",
                start_date: null,
                end_date: null,
                price: ""
            },
            resp: ""
        }
    },
    methods: {
        async submitModalData() {
            try {
                let requestData = {}
                if (this.req.end_date == null) {
                    requestData = {
                        variant_id: parseInt(this.req.variant_id),
                        start_date: this.formatDate(this.req.start_date),
                        price: parseFloat(this.req.price)
                    }

                } else {
                    requestData = {
                        variant_id: parseInt(this.req.variant_id),
                        start_date: this.formatDate(this.req.start_date),
                        end_date: this.formatDate(this.req.end_date),
                        price: parseFloat(this.req.price)
                    }
                }


                const response = await fetch("http://127.0.0.1:9000/product/price", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                this.resp = `Цена продукта успешно добавлена, priceID = ${responseData.Data.price_id}`
            } catch (error) {
                console.error(error)
                this.resp = error
            }
        },

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
        }
    },
}
</script>