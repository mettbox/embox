declare module 'vue-virtual-scroller' {
  import { DefineComponent } from 'vue';

  export interface DynamicScrollerSlot<T = any> {
    item: T;
    index: number;
    active: boolean;
  }

  export const DynamicScroller: DefineComponent<
    {
      items: any[];
      minItemSize: number;
      buffer: number;
      keyField: string;
      emitUpdate: boolean;
    },
    object,
    {
      default: (slotProps: DynamicScrollerSlot) => any;
    }
  >;

  export const DynamicScrollerItem: DefineComponent<Record<string, any>, Record<string, any>, any>;
  export const RecycleScroller: DefineComponent<Record<string, any>, Record<string, any>, any>;

  const VueVirtualScroller: {
    DynamicScroller: typeof DynamicScroller;
    DynamicScrollerItem: typeof DynamicScrollerItem;
    RecycleScroller: typeof RecycleScroller;
  };

  export default VueVirtualScroller;
}
