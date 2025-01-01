package buffer

import "io"

// Read

func Read2[T1 any, T2 any](r io.Reader, t1t Type[T1], t2t Type[T2]) (t1 T1, t2 T2, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	return
}

func Read3[T1 any, T2 any, T3 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3]) (t1 T1, t2 T2, t3 T3, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	return
}

func Read4[T1 any, T2 any, T3 any, T4 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4]) (t1 T1, t2 T2, t3 T3, t4 T4, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	return
}

func Read5[T1 any, T2 any, T3 any, T4 any, T5 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	return
}

func Read6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5], t6t Type[T6]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	if t6, err = t6t.Read(r); err != nil {
		return
	}
	return
}

func Read7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5], t6t Type[T6], t7t Type[T7]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	if t6, err = t6t.Read(r); err != nil {
		return
	}
	if t7, err = t7t.Read(r); err != nil {
		return
	}
	return
}

func Read8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5], t6t Type[T6], t7t Type[T7], t8t Type[T8]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	if t6, err = t6t.Read(r); err != nil {
		return
	}
	if t7, err = t7t.Read(r); err != nil {
		return
	}
	if t8, err = t8t.Read(r); err != nil {
		return
	}
	return
}

func Read9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5], t6t Type[T6], t7t Type[T7], t8t Type[T8], t9t Type[T9]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	if t6, err = t6t.Read(r); err != nil {
		return
	}
	if t7, err = t7t.Read(r); err != nil {
		return
	}
	if t8, err = t8t.Read(r); err != nil {
		return
	}
	if t9, err = t9t.Read(r); err != nil {
		return
	}
	return
}

func Read10[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any](r io.Reader, t1t Type[T1], t2t Type[T2], t3t Type[T3], t4t Type[T4], t5t Type[T5], t6t Type[T6], t7t Type[T7], t8t Type[T8], t9t Type[T9], t10t Type[T10]) (t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10, err error) {
	if t1, err = t1t.Read(r); err != nil {
		return
	}
	if t2, err = t2t.Read(r); err != nil {
		return
	}
	if t3, err = t3t.Read(r); err != nil {
		return
	}
	if t4, err = t4t.Read(r); err != nil {
		return
	}
	if t5, err = t5t.Read(r); err != nil {
		return
	}
	if t6, err = t6t.Read(r); err != nil {
		return
	}
	if t7, err = t7t.Read(r); err != nil {
		return
	}
	if t8, err = t8t.Read(r); err != nil {
		return
	}
	if t9, err = t9t.Read(r); err != nil {
		return
	}
	if t10, err = t10t.Read(r); err != nil {
		return
	}
	return
}

// Write

func Write2[T1 any, T2 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	return
}

func Write3[T1 any, T2 any, T3 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	return
}

func Write4[T1 any, T2 any, T3 any, T4 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	return
}

func Write5[T1 any, T2 any, T3 any, T4 any, T5 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	return
}

func Write6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5, t6t Type[T6], t6 T6) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	if err = t6t.Write(w, t6); err != nil {
		return
	}
	return
}

func Write7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5, t6t Type[T6], t6 T6, t7t Type[T7], t7 T7) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	if err = t6t.Write(w, t6); err != nil {
		return
	}
	if err = t7t.Write(w, t7); err != nil {
		return
	}
	return
}

func Write8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5, t6t Type[T6], t6 T6, t7t Type[T7], t7 T7, t8t Type[T8], t8 T8) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	if err = t6t.Write(w, t6); err != nil {
		return
	}
	if err = t7t.Write(w, t7); err != nil {
		return
	}
	if err = t8t.Write(w, t8); err != nil {
		return
	}
	return
}

func Write9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5, t6t Type[T6], t6 T6, t7t Type[T7], t7 T7, t8t Type[T8], t8 T8, t9t Type[T9], t9 T9) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	if err = t6t.Write(w, t6); err != nil {
		return
	}
	if err = t7t.Write(w, t7); err != nil {
		return
	}
	if err = t8t.Write(w, t8); err != nil {
		return
	}
	if err = t9t.Write(w, t9); err != nil {
		return
	}
	return
}

func Write10[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any, T10 any](w io.Writer, t1t Type[T1], t1 T1, t2t Type[T2], t2 T2, t3t Type[T3], t3 T3, t4t Type[T4], t4 T4, t5t Type[T5], t5 T5, t6t Type[T6], t6 T6, t7t Type[T7], t7 T7, t8t Type[T8], t8 T8, t9t Type[T9], t9 T9, t10t Type[T10], t10 T10) (err error) {
	if err = t1t.Write(w, t1); err != nil {
		return
	}
	if err = t2t.Write(w, t2); err != nil {
		return
	}
	if err = t3t.Write(w, t3); err != nil {
		return
	}
	if err = t4t.Write(w, t4); err != nil {
		return
	}
	if err = t5t.Write(w, t5); err != nil {
		return
	}
	if err = t6t.Write(w, t6); err != nil {
		return
	}
	if err = t7t.Write(w, t7); err != nil {
		return
	}
	if err = t8t.Write(w, t8); err != nil {
		return
	}
	if err = t9t.Write(w, t9); err != nil {
		return
	}
	if err = t10t.Write(w, t10); err != nil {
		return
	}
	return
}
