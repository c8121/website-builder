export default {
    template: `
        <div>
            <asset-select label="URL" v-model="editorContent"></asset-select>
        </div>
    `,

    data() {
        return {
            editorContent: null
        }
    },

    methods: {
        setValue(value) {
            this.editorContent = value
        },

        getValue() {
            return this.editorContent;
        },
    },
}