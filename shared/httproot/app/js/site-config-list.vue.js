import { nextTick } from 'vue'

export default {
    template: `
        <div>
            <table class="table">
                <tbody>
                    <tr v-for="(item, index) in list" role="button">
                        <td>{{ item.Id }}</td>
                        <td><input v-model="item.Name" class="form-control"></td>
                        <td><input v-model="item.Value" class="form-control"></td>
                        <td class="text-end text-nowrap">
                            <button @click.stop="move(index, -1)" v-if="index > 0" class="btn btn-sm m-0 p-0" title="Move up"><i class="bi bi-arrow-up-circle"></i></button>
                            <button @click.stop="move(index, 1)" v-if="index < list.length - 1" class="btn btn-sm m-0 p-0 ms-1" title="Move down"><i class="bi bi-arrow-down-circle"></i></button>
                            <button @click.stop="deleteItem(item)" class="btn btn-sm m-0 p-0 ms-1" title="Delete"><i class="bi bi-trash"></i></button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div class="mt-4">
                <button class="btn btn-primary" @click="saveItems">Save</button>
                <button class="btn btn-secondary ms-2" @click="addItem">Add Config</button>
            </div>
        </div>
    `,

    data() {
        return {
            list: []
        }
    },

    watch: {
        value() {
            this.loadItems();
        }
    },

    methods: {

        loadItems() {
            const self = this;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/sitesconfigs/list', requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.list.splice(0, self.list.length);
                    for (const item of json) {
                        self.list.push(item);
                    }
                });
        },

        reload() {
            this.loadItems();
        },

        move(index, delta) {
            const self = this;

            const target = index + delta;
            if (target < 0 || target >= self.list.length)
                return;

            const ids = self.list.map(item => item.Id);
            [ids[index], ids[target]] = [ids[target], ids[index]];

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({Ids: ids})
            }
            fetch('/sitesconfig/reorder', requestParams)
                .finally(() => self.loadItems());
        },

        deleteItem(item) {
            if(!confirm("Delete " + item.Name + "?")) {
                return;
            }

            const self = this;
            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: "id=" + encodeURI(item.Id)
            }
            fetch('/sitesconfig/delete', requestParams)
                .finally(() => self.loadItems());
        },

        addItem() {
            const self = this;
            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({Name: "name", Value: "value"})
            }
            fetch('/sitesconfig/add', requestParams)
                .finally(() => self.loadItems());
        },

        saveItem(item) {
            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(item)
            }
            return fetch('/sitesconfig/add', requestParams)
        },

        saveItems() {
            const self = this;
            const fetches = [];
            for( const item of this.list) {
                fetches.push(this.saveItem(item));
            }
            Promise.all(fetches)
                .finally(() => self.loadItems());
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadItems();
    }
}