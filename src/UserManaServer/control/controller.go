package control

// func AddUser(arg_info map[string][string]) map[string][string] {
// 	return nil
// }

// func DelUser(arg_info map[string][string]) map[string][string] {
// 	return nil
// }

// func DispatchAuth(arg_info map[string][string]) map[string][string]{
// 	return nil
// }

// func ModifyUserMsg(arg_info map[string][string]) map[string][string] {
// 	return nil
// }

// func AddBillboard(arg_info map[string][string]) {
// 	return nil
// }

// func ViewBillboard(arg_info map[string][string]) map[string][string] {
// 	return nil
// }

// func ModifyBillboard(arg_info map[string][string]) map[string][string] {
// 	return nil
// }
// func DelBillboard(arg_info map[string][string]) map[string][string] {
// 	return nil
// }
// func ViewMsg(arg_info map[string][string]) map[string][string] {
// 	return nil
// }
// func ModifyMsg(arg_info map[string][string]) map[string][string] {
// 	return nil
// }
var G_op_func map[string]func(map[string]string) map[string]string = make(map[string]func(map[string]string) map[string]string)
