import { nextTick } from 'vue'

export default {

    name: 'page-tree',

    template: `
        <div>
            <div v-for="(item, index) in list" class="text-nowrap text-start">
                
                <div class="hover-container">
                    <button @click="item.showChildren = !item.showChildren" class="btn m-0 p-0">
                        <span v-if="!item.showChildren"><i class="bi bi-folder-plus"></i></span>
                        <span v-else><i class="bi bi-folder-minus"></i></span>
                    </button>
                    <button @click="item.showChildren = true; itemClicked(item)" class="btn m-0 p-0">
                        <span v-if="!item.Name">({{ item.Id }})</span>{{ item.Name }}
                    </button>
                    <span v-if="sortable" class="ms-2 hover-content">
                        <button @click="move(index, -1)" v-if="index > 0" class="btn m-0 p-0" title="Move up"><i class="bi bi-arrow-up-circle"></i></button>
                        <button @click="move(index, 1)" v-if="index < list.length - 1" class="btn m-0 p-0" title="Move down"><i class="bi bi-arrow-down-circle"></i></button>
                    </span>
                </div>
                
                <div class="ms-2">
                <page-tree ref="childTrees" v-if="item.showChildren" :value="item.Id" :sortable="sortable"
                        @componentEvent="childClicked"></page-tree>
                </div>
            </div>
        </div>
    `,

    props: {
        value: "",
        sortable: {
            type: Boolean,
            default: false
        }
    },

    data() {
        return {
            parentId: 0,
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

            const showChildren = [];
            for (const item of self.list) {
                if (item.showChildren)
                    showChildren.push(item.Id);
            }

            const requestOptions = {
                method: 'GET'
            }
            const parentId = self.value ? self.value : 0;
            fetch('/pages/list/' + parentId, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.list.splice(0, self.list.length);
                    for (const item of json) {
                        item.showChildren = false;
                        self.list.push(item);
                    }

                    //Update showChildren after populating list to trigger reload
                    nextTick(() => {
                        for (const item of self.list)
                            item.showChildren = showChildren.includes(item.Id);
                    });
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
                body: JSON.stringify({ParentId: self.value ? self.value : 0, Ids: ids})
            }
            fetch('/pages/reorder', requestParams)
                .finally(() => self.loadItems());
        },

        itemClicked(item) {
            this.$emit('componentEvent', 'pageClick', 'page-tree', item);
        },

        childClicked() {
            this.$emit('componentEvent', 'pageClick', 'page-tree', arguments[2]);
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadItems();
    }
}