INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (1007, 6, '充提币渠道管理', '', 1, NULL, 'WalletChain', 'Index', 0, 0, '2024-07-13 12:51:12.094', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1007, '充提币渠道编辑', '', 0, NULL, 'WalletChain', 'Edit', 0, 0, '2024-07-13 12:51:12.112', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1007, '充提币渠道删除', '', 0, NULL, 'WalletChain', 'Delete', 0, 0, '2024-07-13 12:51:12.114', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1007, '充提币渠道添加', '', 0, NULL, 'WalletChain', 'Add', 0, 0, '2024-07-13 12:51:12.115', '2024-06-09 12:03:11.987', NULL);
INSERT INTO `nov_role` (`id`, `pid`, `name`, `icon`, `is_menu`, `desc`, `module`, `action`, `sort`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES (NULL, 1007, '充提币渠道批量删除', '', 0, NULL, 'WalletChain', 'Deletebatch', 0, 0, '2024-07-13 12:51:12.117', '2024-06-09 12:03:11.987', NULL);
package model

import "time" 
type WalletChain struct {
Id int `comment:"" types:"" text:"" json:"id" form:"id" range:"" edit:""  `
Title string `comment:"" types:"" text:"" json:"title" form:"title" range:"" edit:""  `
Icon string `comment:"" types:"" text:"" json:"icon" form:"icon" range:"" edit:""  `
Pid int `comment:"" types:"" text:"" json:"pid" form:"pid" range:"" edit:""  `
Chain string `comment:"" types:"" text:"" json:"chain" form:"chain" range:"" edit:""  `
Label string `comment:"" types:"" text:"" json:"label" form:"label" range:"" edit:""  `
CreateTime time.Time `comment:"" types:"" text:"" json:"create_time" form:"create_time" range:"" edit:""  `
UpdateTime time.Time `comment:"" types:"" text:"" json:"update_time" form:"update_time" range:"" edit:""  `
Remarks string `comment:"" types:"" text:"" json:"remarks" form:"remarks" range:"" edit:""  `
Statuswithdraw int `comment:"" types:"" text:"" json:"statuswithdraw" form:"statuswithdraw" range:"" edit:""  `
Statusrecharge int `comment:"" types:"" text:"" json:"statusrecharge" form:"statusrecharge" range:"" edit:""  `
Address string `comment:"" types:"" text:"" json:"address" form:"address" range:"" edit:""  `
IconUrl string `comment:"" types:"" text:"" json:"icon_url" form:"icon_url" range:"" edit:""  `
SubCategory []*model.WalletChain `comment:"" types:"" text:"" json:"sub_category" form:"sub_category" range:"" edit:""  `
}