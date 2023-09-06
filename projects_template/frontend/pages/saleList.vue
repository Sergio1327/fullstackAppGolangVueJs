<template>
    <section>
        <h1 class="label is-text-center mb-5">Продажи</h1>

        <div class="is-flex is-align-items-center is-justify-content-flex-end">
            <button class="button py-3 px-5 is-warning" type="submit" @click="openModal">Добавить продажу</button>
            <AddSaleForm v-if="modalVisible" :stockOptions="stockList" :variantOptions="variantIDs"
                :modalVisible="modalVisible" @closeModal="closeModal" />

        </div>

        <div>
            <div class="field">
                <label class="label">Дата начала продаж</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="req.startDate">
                </div>
            </div>

            <div class="field">
                <label class="label">Дата конца продаж</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="req.endDate">
                </div>
            </div>

            <div class="field">
                <label class="label">Лимит вывода</label>
                <div class="control">
                    <div class="select">
                        <select v-model="req.limit">
                            <option value="1">1</option>
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="50">50</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label">Название продукта</label>
                <div class="control">
                    <input class="input" type="text" v-model="req.productName">
                </div>
            </div>

            <div class="field">
                <label class="label">ID склада</label>
                <div class="control">
                    <input class="input" type="number" v-model="req.storageId">
                </div>
            </div>

            <div class="field is-grouped">
                <div class="control">
                    <button class="button is-link" @click="sendRequest">Найти продажи</button>
                </div>
            </div>
        </div>


        <b-table class="mt-6" :data="saleListData" :hoverable="isHoverable" :striped="isStriped"
            :columns="columns"></b-table>
    </section>
</template>



<script>
import AddSaleForm from "~/components/AddSaleForm.vue";
export default {
    components: {
        AddSaleForm
    }, methods: {
        async openModal() {
            try {

                const response = await fetch("http://127.0.0.1:9000/stock_list", {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    },

                })
                const responseData = await response.json()
                const data = responseData.Data.stock_list
                console.log(data)
                this.stockList = data.map(e => {
                    return {
                        Option: e.StorageID,
                        Value: e.StorageID,
                        StorageName: e.StorageName
                    }
                })

                const response2 = await fetch("http://127.0.0.1:9000/product_list?limit=999999", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json"
                    }
                })

                const responseData2 = await response2.json()

                responseData2.Data.product_list.forEach(product => {
                    product.VariantList.forEach(variant => {
                        const newItem = {
                            ProductName: product.Name,
                            Option: variant.variant_id,
                            Weight: variant.weight,
                            Unit: variant.unit
                        };
                        this.variantIDs.push(newItem);
                    });
                });

                this.variantIDs.sort((a, b) => a.Option - b.Option)
                this.modalVisible = true;

            } catch (error) {
                console.error(error)
            }
        },

        closeModal() {
            this.variantIDs = []
            this.modalVisible = false;
        },

        handleData(data) {
            this.saleListData = data
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

        async LoadSales() {
            try {
                const formattedStartDate = this.formatDate(this.startDate);
                const formattedEndDate = this.formatDate(this.endDate);

                const requestData = {
                    start_date: formattedStartDate,
                    end_date: formattedEndDate,
                    limit: parseInt(this.limit),
                    product_name: this.productName,
                    storage_id: parseInt(this.storageId)
                };
                const response = await fetch("http://127.0.0.1:9000/sales", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData)
                })

                const responseData = await response.json()
                this.saleListData = responseData.Data.sale_list

            } catch (error) {
                console.error(error)
            }
        },
        async sendRequest() {

            const formattedStartDate = this.formatDate(this.req.startDate);
            const formattedEndDate = this.formatDate(this.req.endDate);

            const requestData = {
                start_date: formattedStartDate,
                end_date: formattedEndDate,
                limit: +this.req.limit,
                product_name: this.req.productName,
                storage_id: +this.req.storageId
            };

            console.log(requestData)
            try {
                const response = await fetch("http://127.0.0.1:9000/sales", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                });

                const responseData = await response.json();
                console.log(responseData)
                const data = responseData.Data.sale_list
                this.saleListData = data

            } catch (error) {
                console.error("Ошибка при отправке запроса:", error);
                console.log(error)
            }
        },

    }, data() {
        return {
            req: {
                startDate: null,
                endDate: null,
                limit: 1,
                productName: "",
                storageId: 1,
            },

            modalVisible: false,
            stockList: [],
            variantIDs: [],

            startDate: new Date("2020-08-20T00:00:00"),
            endDate: new Date("2029-08-20T00:00:00"),
            limit: 2,
            productName: "",
            storageId: 1,

            saleListData: [],
            isStriped: true,
            isHoverable: true,

            columns: [
                {
                    field: "SaleID",
                    label: "ID Продажи",
                },
                {
                    field: "ProductName",
                    label: "Название продукта"
                },
                {
                    field: "variant_id",
                    label: "ID варианта"
                },
                {
                    field: "storage_id",
                    label: "ID склада"
                },
                {
                    field: "SoldAt",
                    label: "Дата продажи"
                },
                {
                    field: "quantity",
                    label: "Колличество"
                },
                {
                    field: "TotalPrice",
                    label: "Общая цена продажи"
                }

            ]
        }
    },
    async mounted() {
        await this.LoadSales()
    }
}
</script>

