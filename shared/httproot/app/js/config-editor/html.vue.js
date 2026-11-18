export default {
    template: `
        <div>
            <template v-if="this.editorContent">
                <neiki-editor v-model="editorContent"></neiki-editor>
            </template>
        </div>
    `,

    data() {
        return {
            editorContent: ""
        }
    },

    methods: {
        setValue(value) {
            this.editorContent = value ? value : "&nbsp;"; //TODO find why editor not showing if no value
        },

        getValue() {
            return this.editorContent;
        },
    },
}