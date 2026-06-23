INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (1026, 54, '产品管理', '', 1, NULL, 'Goods', 'Index', 0, 0, '2024-07-13 12:51:12.094', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1026, '产品编辑', '', 0, NULL, 'Goods', 'Edit', 0, 0, '2024-07-13 12:51:12.112', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1026, '产品删除', '', 0, NULL, 'Goods', 'Delete', 0, 0, '2024-07-13 12:51:12.114', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1026, '产品添加', '', 0, NULL, 'Goods', 'Add', 0, 0, '2024-07-13 12:51:12.115', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1026, '产品批量删除', '', 0, NULL, 'Goods', 'Deletebatch', 0, 0, '2024-07-13 12:51:12.117', '2024-06-09 12:03:11.987', NULL);
package model

import "time" 
type Goods struct {
Model model.Model `comment:"" types:"" text:"" json:"model" form:"model" range:"" edit:""  `
GoodsCate int `comment:"商品类别" types:"" text:"" json:"goods_cate" form:"goods_cate" range:"" edit:""  `
GoodsName string `comment:"商品类别" types:"" text:"" json:"goods_name" form:"goods_name" range:"" edit:""  `
GoodsProperty string `comment:"" types:"" text:"" json:"goods_property" form:"goods_property" range:"" edit:"0"  `
GoodsDesc string `comment:"商品描述" types:"" text:"" json:"goods_desc" form:"goods_desc" range:"" edit:""  `
GoodsContent string `comment:"商品信息" types:"" text:"" json:"goods_content" form:"goods_content" range:"" edit:""  `
UnitPrice float64 `comment:"商品单价" types:"" text:"" json:"unit_price" form:"unit_price" range:"" edit:""  `
FavorablePrice float64 `comment:"优惠价格" types:"" text:"" json:"favorable_price" form:"favorable_price" range:"" edit:""  `
GoodsStock uint64 `comment:"商品库存" types:"" text:"" json:"goods_stock" form:"goods_stock" range:"" edit:""  `
GoodsCover string `comment:"商品封面图" types:"" text:"" json:"goods_cover" form:"goods_cover" range:"" edit:""  `
GoodsSlides string `comment:"" types:"" text:"" json:"goods_slides" form:"goods_slides" range:"" edit:""  `
GoodsStatus uint64 `comment:"状态" types:"" text:"" json:"goods_status" form:"goods_status" range:"" edit:""  `
}