<template>
    <section>
        <h1 class="label is-text-center mb-5">Продажи</h1>

        <div class="is-flex is-align-items-center is-justify-content-flex-end">
            <button class="button py-3 px-5 is-warning" type="submit" @click="openModal">Добавить продажу</button>
            <AddSaleForm v-if="modalVisible" :stockOptions="stockList" :variantOptions="variantIDs"
                :modalVisible="modalVisible" @closeModal="closeModal" />
        </div>
        <div class="">
            <div class="field">
                <label class="label">Дата начала продаж</label>
                <div class="control">
                    <b-datepicker class="form" v-model="req.startDate" ref="datepicker" expanded
                        placeholder="Введите дату начала продаж">
                    </b-datepicker>
                </div>
            </div>

            <div class="field">
                <label class="label">Дата конца продаж</label>
                <div class="control">
                    <b-datepicker class="form" v-model="req.endDate" ref="datepicker" expanded
                        placeholder="Введите дату конца продаж">
                    </b-datepicker>
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
                    <input class="input form" placeholder="Введите название продукта" type="text" v-model="req.productName">
                </div>
            </div>

            <b-field label="Выберите склад">
                <b-select v-model="req.storageId" type="text" placeholder="Выберите склад" required>
                    <option v-for="s in stockOptions" :value="s.Option" :key="s.Value">{{ s.StorageName }}</option>
                </b-select>
            </b-field>

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
    },

    methods: {
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
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        },

        async LoadSales() {
            try {
                const requestData = {
                    start_date: this.startDate,
                    end_date: this.endDate,
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
                this.formateDate(this.saleListData)
            } catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        },

        async sendRequest() {
            const requestData = {
                start_date: this.req.startDate,
                end_date: this.req.endDate,
                limit: +this.req.limit,
                product_name: this.req.productName,
                storage_id: +this.req.storageId
            };

            try {
                const response = await fetch("http://127.0.0.1:9000/sales", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                });

                const responseData = await response.json();
                const data = responseData.Data.sale_list
                this.saleListData = data
                this.formateDate(this.saleListData)

            } catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error("Ошибка при отправке запроса:", error);
            }
        },

        async fetchStockList() {
            try {
                const response = await fetch('http://localhost:9000/stock_list', {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                const responseData = await response.json();
                const data = responseData.Data.stock_list

                this.stockOptions = data.map(e => {
                    return {
                        Option: e.StorageID,
                        Value: e.StorageID,
                        StorageName: e.StorageName
                    }
                })
            } catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error('Error fetching warehouses:', error);
            }
        },

        closeModal() {
            this.variantIDs = []
            this.modalVisible = false;
        },

        formateDate(data) {
            const userLocale = "ru-RU"; // Укажите нужную локаль пользователя
            const options = {
                year: "numeric",
                month: "long",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
                second: "2-digit",
                timeZoneName: "short",
            };
            data.forEach(sale => {
                const soldAtDate = new Date(sale.SoldAt)
                sale.SoldAt = soldAtDate.toLocaleString(userLocale, options)
            })
        }
    },
    data() {
        return {
            date: null,
            req: {
                startDate: null,
                endDate: null,
                limit: 1,
                productName: "",
                storageId: 1,
            },

            stockOptions: [],
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

            ],
        }
    },
    async mounted() {
        await this.LoadSales()
        await this.fetchStockList()
    }
}
</script>

<style>
.sales_form {
    max-width: 1000px !important;
}

.form {
    width: 30% !important;
}

@media (max-width:1000px) {
    .form {
        width: 100% !important;
    }
}
</style>

