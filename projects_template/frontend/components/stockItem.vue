<template>
    <li>
        <div class="card ">
            <div class="card-content is-flex is-align-items-center is-justify-content-space-between ">
                <div class="content column is-text-center">
                    <div>ID склада : <strong> {{ stock.StorageID }}</strong></div>
                    <div>Название склада : <strong> {{ stock.StorageName }}</strong></div>
                </div>
                <div class="image is-32x32 mx-4">
                    <img class="is-img trash" src="@/assets/trash.png" alt="">
                </div>
                <div class="image is-128x128 is-flex is-align-items-center">
                    <img src="@/assets/stockList.jpg" class="is-img" alt="">
                </div>
                <div>{{ msg }}</div>
            </div>
        </div>
    </li>
</template>


<script>
export default {
    props: {
        stock: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            msg: ""
        }
    }
    ,

    methods: {
        async deleteStock() {
            const requestData = {
                storage_id: this.stock.StorageID,
                storage_name: this.stock.StorageName
            }
            try {
                const response = await fetch("http://127.0.0.1:9000/stock/delete", {
                    method: "DELETE",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                if (response.ok) {
                    this.msg = "склад успешно удален"

                    setTimeout(() => {
                        this.msg = ""
                    }, 2000)

                } else {
                    this.msg = response
                }
            } catch (error) {
                console.log(error)
                this.msg = error
            }

        }
    }
}
</script>
<style scoped >
content div strong {
    color: #000;
    font-weight: 600;
}

.trash {
    transition: 200ms linear;
}

.trash:hover {
    transform: translateY(-3px) !important;
}

.card {
    transition: 200ms linear;
    box-shadow: 1px 1px 5px rgb(73, 49, 5);
}

.card:hover {
    box-shadow: 2px 2px 10px rgb(255, 183, 49);
    transform: translateY(-5px);
}
</style>