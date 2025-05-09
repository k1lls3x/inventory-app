export namespace dashboard {
	
	export class DashboardData {
	    total_stock: number;
	    item_count: number;
	    monthly_orders: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_stock = source["total_stock"];
	        this.item_count = source["item_count"];
	        this.monthly_orders = source["monthly_orders"];
	    }
	}

}

