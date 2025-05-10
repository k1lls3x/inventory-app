export namespace dashboard {
	
	export class DashboardData {
	    total_stock: number;
	    item_count: number;
	    monthly_orders: number;
	    new_items: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_stock = source["total_stock"];
	        this.item_count = source["item_count"];
	        this.monthly_orders = source["monthly_orders"];
	        this.new_items = source["new_items"];
	    }
	}

}

export namespace model {
	
	export class ItemTurnoverByWarehouse {
	    warehouse: string;
	    name: string;
	    sku: string;
	    received: number;
	    shipped: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemTurnoverByWarehouse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.warehouse = source["warehouse"];
	        this.name = source["name"];
	        this.sku = source["sku"];
	        this.received = source["received"];
	        this.shipped = source["shipped"];
	    }
	}
	export class ItemWithStock {
	    name: string;
	    sku: string;
	    warehouse: string;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemWithStock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.sku = source["sku"];
	        this.warehouse = source["warehouse"];
	        this.quantity = source["quantity"];
	    }
	}
	export class Stock {
	    stock_id: number;
	    item_id: number;
	    quantity: number;
	    warehouse_id: number;
	
	    static createFrom(source: any = {}) {
	        return new Stock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stock_id = source["stock_id"];
	        this.item_id = source["item_id"];
	        this.quantity = source["quantity"];
	        this.warehouse_id = source["warehouse_id"];
	    }
	}

}

