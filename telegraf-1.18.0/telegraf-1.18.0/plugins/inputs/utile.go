package inputs


import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "time"
    "strings"
)



func  Check_swith(Swith_data_type []string,field map[string]interface{}){
      if len(Swith_data_type)>0{
         dmp:=make(map[string]string)
         for _,_v:=range Swith_data_type{
             x:=strings.Split(_v,"@")
             dmp[x[0]]=x[1]
         }    
         for k,v:=range dmp{ 
             if v1,ok:=field[k];ok{
                //fmt.Printf("cached this : %s  value is %v  type is : %v\n",k,v1,reflect.TypeOf(v1))
                switch v{
                   case "int":
                       if (reflect.TypeOf(v1).Name()=="int64") {
                           return
                       }    
                       var dst uint64 
                       convertAssign(&dst, v1)
                       field[k]=dst
                       // fmt.Printf("after switch this : %s  value is %v  type is : %v\n",k,dst,reflect.TypeOf(v1))
                   case "float":
                       if (reflect.TypeOf(v1).Name()=="float64") {
                           return
                       }    
                       var dst float64 
                       convertAssign(&dst, v1)
                       field[k]=dst
                       // fmt.Printf("after switch this : %s  value is %v  type is : %v\n",k,dst,reflect.TypeOf(v1))
                  case "string":
                       if (reflect.TypeOf(v1).Name()=="string") {
                           return
                       }    
                       var dst string 
                       convertAssign(&dst, v1)
                       field[k]=dst
                }    
             }    
         }    
     }    
} 

//  add conver function

 
// RawBytes is a byte slice that holds a reference to memory owned by
// the database itself. After a Scan into a RawBytes, the slice is only
// valid until the next call to Next, Scan, or Close.
type RawBytes []byte
 
var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error
 
// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest, src interface{}) error {
    // Common cases, without reflect.
    switch s := src.(type) {
    case string:
        switch d := dest.(type) {
        case *string:
            if d == nil {
                return errNilPtr
            }
            *d = s
            return nil
        case *[]byte:
            if d == nil {
                return errNilPtr
            }
            *d = []byte(s)
            return nil
        case *RawBytes:
            if d == nil {
                return errNilPtr
            }
            *d = append((*d)[:0], s...)
            return nil
        }
    case []byte:
        switch d := dest.(type) {
        case *string:
            if d == nil {
                return errNilPtr
            }
            *d = string(s)
            return nil
        case *interface{}:
            if d == nil {
                return errNilPtr
            }
            *d = cloneBytes(s)
            return nil
        case *[]byte:
            if d == nil {
                return errNilPtr
            }
            *d = cloneBytes(s)
            return nil
        case *RawBytes:
            if d == nil {
                return errNilPtr
            }
            *d = s
            return nil
        }
    case time.Time:
        switch d := dest.(type) {
        case *time.Time:
            *d = s
            return nil
        case *string:
            *d = s.Format(time.RFC3339Nano)
            return nil
        case *[]byte:
            if d == nil {
                return errNilPtr
            }
            *d = []byte(s.Format(time.RFC3339Nano))
            return nil
        case *RawBytes:
            if d == nil {
                return errNilPtr
            }
            *d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
            return nil
        }
    case nil:
        switch d := dest.(type) {
        case *interface{}:
            if d == nil {
                return errNilPtr
            }
            *d = nil
            return nil
        case *[]byte:
            if d == nil {
                return errNilPtr
            }
            *d = nil
            return nil
        case *RawBytes:
            if d == nil {
                return errNilPtr
            }
            *d = nil
            return nil
        }
    }
 
    var sv reflect.Value
 
    switch d := dest.(type) {
    case *string:
        sv = reflect.ValueOf(src)
        switch sv.Kind() {
        case reflect.Bool,
            reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
            reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
            reflect.Float32, reflect.Float64:
            *d = asString(src)
            return nil
        }
    case *[]byte:
        sv = reflect.ValueOf(src)
        if b, ok := asBytes(nil, sv); ok {
            *d = b
            return nil
        }
    case *RawBytes:
        sv = reflect.ValueOf(src)
        if b, ok := asBytes([]byte(*d)[:0], sv); ok {
            *d = RawBytes(b)
            return nil
        }
    case *bool:
        bv, err := Bool.ConvertValue(src)
        if err == nil {
            *d = bv.(bool)
        }
        return err
    case *interface{}:
        *d = src
        return nil
    }
 
    if scanner, ok := dest.(Scanner); ok {
        return scanner.Scan(src)
    }
 
    dpv := reflect.ValueOf(dest)
    if dpv.Kind() != reflect.Ptr {
        return errors.New("destination not a pointer")
    }
    if dpv.IsNil() {
        return errNilPtr
    }
 
    if !sv.IsValid() {
        sv = reflect.ValueOf(src)
    }
 
    dv := reflect.Indirect(dpv)
    if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
        switch b := src.(type) {
        case []byte:
            dv.Set(reflect.ValueOf(cloneBytes(b)))
        default:
            dv.Set(sv)
        }
        return nil
    }
 
    if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
        dv.Set(sv.Convert(dv.Type()))
        return nil
    }
 
    // The following conversions use a string value as an intermediate representation
    // to convert between various numeric types.
    //
    // This also allows scanning into user defined types such as "type Int int64".
    // For symmetry, also check for string destination types.
    switch dv.Kind() {
    case reflect.Ptr:
        if src == nil {
            dv.Set(reflect.Zero(dv.Type()))
            return nil
        }
        dv.Set(reflect.New(dv.Type().Elem()))
        return convertAssign(dv.Interface(), src)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        s := asString(src)
        i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
        if err != nil {
            err = strconvErr(err)
            return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
        }
        dv.SetInt(i64)
        return nil
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        s := asString(src)
        u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
        if err != nil {
            err = strconvErr(err)
            return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
        }
        dv.SetUint(u64)
        return nil
    case reflect.Float32, reflect.Float64:
        s := asString(src)
        f64, err := strconv.ParseFloat(s, dv.Type().Bits())
        if err != nil {
            err = strconvErr(err)
            return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
        }
        dv.SetFloat(f64)
        return nil
    case reflect.String:
        switch v := src.(type) {
        case string:
            dv.SetString(v)
            return nil
        case []byte:
            dv.SetString(string(v))
            return nil
        }
    }
 
    return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}
 
func strconvErr(err error) error {
    if ne, ok := err.(*strconv.NumError); ok {
        return ne.Err
    }
    return err
}
 
func cloneBytes(b []byte) []byte {
    if b == nil {
        return nil
    }
    c := make([]byte, len(b))
    copy(c, b)
    return c
}
 
func asString(src interface{}) string {
    switch v := src.(type) {
    case string:
        return v
    case []byte:
        return string(v)
    }
    rv := reflect.ValueOf(src)
    switch rv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return strconv.FormatInt(rv.Int(), 10)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return strconv.FormatUint(rv.Uint(), 10)
    case reflect.Float64:
        return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
    case reflect.Float32:
        return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
    case reflect.Bool:
        return strconv.FormatBool(rv.Bool())
    }
    return fmt.Sprintf("%v", src)
}
 
func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
    switch rv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return strconv.AppendInt(buf, rv.Int(), 10), true
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return strconv.AppendUint(buf, rv.Uint(), 10), true
    case reflect.Float32:
        return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
    case reflect.Float64:
        return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
    case reflect.Bool:
        return strconv.AppendBool(buf, rv.Bool()), true
    case reflect.String:
        s := rv.String()
        return append(buf, s...), true
    }
    return
}
 
// Value is a value that drivers must be able to handle.
// It is either nil, a type handled by a database driver's NamedValueChecker
// interface, or an instance of one of these types:
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
type Value interface{}
 
type boolType struct{}
var Bool boolType
func (boolType) String() string { return "Bool" }
func (boolType) ConvertValue(src interface{}) (Value, error) {
    switch s := src.(type) {
    case bool:
        return s, nil
    case string:
        b, err := strconv.ParseBool(s)
        if err != nil {
            return nil, fmt.Errorf("sql/driver: couldn't convert %q into type bool", s)
        }
        return b, nil
    case []byte:
        b, err := strconv.ParseBool(string(s))
        if err != nil {
            return nil, fmt.Errorf("sql/driver: couldn't convert %q into type bool", s)
        }
        return b, nil
    }
 
    sv := reflect.ValueOf(src)
    switch sv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        iv := sv.Int()
        if iv == 1 || iv == 0 {
            return iv == 1, nil
        }
        return nil, fmt.Errorf("sql/driver: couldn't convert %d into type bool", iv)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        uv := sv.Uint()
        if uv == 1 || uv == 0 {
            return uv == 1, nil
        }
        return nil, fmt.Errorf("sql/driver: couldn't convert %d into type bool", uv)
    }
 
    return nil, fmt.Errorf("sql/driver: couldn't convert %v (%T) into type bool", src, src)
}
 
type Scanner interface {
    // Scan assigns a value from a database driver.
    //
    // The src value will be of one of the following types:
    //
    //    int64
    //    float64
    //    bool
    //    []byte
    //    string
    //    time.Time
    //    nil - for NULL values
    //
    // An error should be returned if the value cannot be stored
    // without loss of information.
    //
    // Reference types such as []byte are only valid until the next call to Scan
    // and should not be retained. Their underlying memory is owned by the driver.
    // If retention is necessary, copy their values before the next call to Scan.
    Scan(src interface{}) error
}
